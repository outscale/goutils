package sanitize

import (
	"bytes"
	"io"
	"net/http"
	"slices"
	"strings"
)

const (
	Request = iota
	Response
)

// Request/Response - service - attribute - sensitivity
type serviceFieldMap map[int]map[string]map[string]string

func (s serviceFieldMap) sensivityFor(typ int, service, attribute string) string {
	if s[typ][service] == nil {
		return ""
	}
	return s[typ][service][attribute]
}

//go:generate go run generate/main.go -package github.com/outscale/osc-sdk-go/v3/pkg/osc --suffix Request --output http_oapi_request.go --service oapi --type Request
//go:generate go run generate/main.go -package github.com/outscale/osc-sdk-go/v3/pkg/osc --suffix Response --output http_oapi_response.go --service oapi --type Response
//go:generate go run generate/main.go -package github.com/outscale/osc-sdk-go/v3/pkg/oks --suffix Request --output http_oks_request.go --service oks --type Request
//go:generate go run generate/main.go -package github.com/outscale/osc-sdk-go/v3/pkg/oks --suffix Response --output http_oks_response.go --service oks --type Response
var serviceFields = serviceFieldMap{
	Request:  {},
	Response: {},
}

type HTTPSanitizer struct {
	sensitivities []string
	sanitize      StringSanitizer
	redact        redactFunc
}

type HTTPOption func(*HTTPSanitizer)

// RedactHTTP configures the redaction function.
func RedactHTTP(fn redactFunc) HTTPOption {
	return func(s *HTTPSanitizer) {
		s.redact = fn
	}
}

// MatchHTTP matches the sensitivities that need to be redacted.
func MatchHTTP(sensitivity ...string) HTTPOption {
	return func(s *HTTPSanitizer) {
		s.sensitivities = append(s.sensitivities, sensitivity...)
	}
}

// NewHTTP will create a sanitizer that will redact:
// * all fields marked as pii or sensitive from a request/response body,
// * headers containing SK.
func NewHTTP(opts ...HTTPOption) HTTPSanitizer {
	if len(opts) == 0 {
		opts = []HTTPOption{MatchHTTP(PII, Sensitive)}
	}
	s := HTTPSanitizer{
		redact: RedactAll,
	}
	for _, opt := range opts {
		opt(&s)
	}
	s.sanitize = New(MatchRegexps(SecretKey), RedactString(s.redact))
	return s
}

func getService(h http.Header) string {
	// AWS4-HMAC-SHA256 Credential=[ak]/20260908/[region]/osc/aws4_request, SignedHeaders=host;x-amz-date, Signature=[...]
	for parts := range strings.SplitSeq(h.Get("Authorization"), " ") {
		kv := strings.TrimRight(parts, ",")
		key, value, _ := strings.Cut(kv, "=")
		if key == "Credential" {
			creds := strings.Split(value, "/")
			switch {
			case len(creds) < 4 || creds[3] == "osc":
				return "oapi"
			default:
				return creds[3]
			}
		}
	}
	return "oks"
}

// SanitizeRequest will sanitize a http.Request.
// A copy of the source value is returned with redacted headers/fields.
func (s HTTPSanitizer) SanitizeRequest(req *http.Request) *http.Request {
	req = req.Clone(req.Context())

	s.sanitizeMapSlice(req.Header)
	s.sanitizeMapSlice(req.Form)
	s.sanitizeMapSlice(req.PostForm)
	if req.MultipartForm != nil {
		s.sanitizeMapSlice(req.MultipartForm.Value)
	}
	// TODO: sanitize query

	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return req
		}
		buf, err := io.ReadAll(body)
		if err != nil {
			req.Body = nil
			return req
		}
		buf = s.sanitizeJSON(buf, Request, getService(req.Header))
		req.Body = io.NopCloser(bytes.NewBuffer(buf))
	}
	return req
}

// SanitizeRequest will sanitize a http.Response.
// A copy of the source value is returned with redacted headers/fields.
func (s HTTPSanitizer) SanitizeResponse(resp *http.Response) *http.Response {
	nresp := &http.Response{}
	*nresp = *resp
	nresp.Header = resp.Header.Clone()
	s.sanitizeMapSlice(resp.Header)

	buf, err := io.ReadAll(resp.Body)
	if err != nil || resp.Request == nil {
		_ = resp.Body.Close()
		resp.Body = nil
		return resp
	}
	_ = resp.Body.Close()

	resp.Body = io.NopCloser(bytes.NewBuffer(buf))

	buf = s.sanitizeJSON(buf, Response, getService(resp.Request.Header))
	nresp.Body = io.NopCloser(bytes.NewBuffer(buf))

	return nresp
}

func (s HTTPSanitizer) sanitizeMapSlice[T ~map[string][]string](m T) {
	if m == nil {
		return
	}
	for k, vs := range m {
		for i, v := range vs {
			m[k][i] = s.sanitize.Sanitize(v)
		}
	}
}

const (
	tokenBufferSize = 128
)

const (
	stateNone = iota
	stateObject
	stateKey
	stateExpectColon
	stateExpectValue
	stateValue
)

type jsonState []int

func (s jsonState) current() int {
	if len(s) == 0 {
		return stateNone
	}
	return s[len(s)-1]
}

func (s *jsonState) pop() {
	if len(*s) <= 1 {
		return
	}
	*s = (*s)[:len(*s)-1]
}

func (s *jsonState) add(state int) {
	*s = append(*s, state)
}

func (s HTTPSanitizer) sanitizeJSON(src []byte, typ int, service string) []byte {
	dst := make([]byte, 0, len(src))
	state := jsonState{}
	token := make([]byte, 0, tokenBufferSize)
	redact := false
	escape := false
	for _, c := range src {
		switch {
		case c == '{':
			if state.current() == stateExpectValue {
				state.pop()
			}
			state.add(stateObject)
		case state.current() == stateNone:
			// not an object, abort
			return src
		case c == '}':
			state.pop()
		case c == '"' && !escape:
			switch state.current() {
			case stateObject:
				state.add(stateKey)
			case stateKey:
				redact = slices.Contains(s.sensitivities, serviceFields.sensivityFor(typ, service, string(token)))

				token = token[:0]
				state.pop()
				state.add(stateExpectColon)
			case stateExpectValue:
				state.pop()
				state.add(stateValue)
			case stateValue:
				if redact {
					sanitized := s.redact(string(token))
					dst = dst[:len(dst)-len(token)]
					dst = append(dst, []byte(sanitized)...)
				}

				token = token[:0]
				state.pop()
			}
		case c == ':' && state.current() == stateExpectColon:
			state.pop()
			state.add(stateExpectValue)
		case state.current() == stateKey || state.current() == stateValue:
			token = append(token, c)
		}
		if c == '\\' {
			escape = !escape
		} else {
			escape = false
		}

		dst = append(dst, c)
	}
	return dst
}

var defaultHTTP = NewHTTP()

func HTTPRequest(req *http.Request) *http.Request {
	return defaultHTTP.SanitizeRequest(req)
}

func HTTPResponse(resp *http.Response) *http.Response {
	return defaultHTTP.SanitizeResponse(resp)
}
