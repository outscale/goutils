package sanitize_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/outscale/goutils/sdk/sanitize"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testLogger struct {
	assertRequest      func(req *http.Request)
	assertRequestBody  func(body string)
	assertResponse     func(resp *http.Response)
	assertResponseBody func(body string)
}

func (l testLogger) Request(ctx context.Context, req any)   {}
func (l testLogger) Response(ctx context.Context, resp any) {}
func (l testLogger) RequestHttp(ctx context.Context, req *http.Request) {
	req = sanitize.HTTPRequest(req)
	if l.assertRequest != nil {
		l.assertRequest(req)
	}
	if l.assertRequestBody != nil {
		body, _ := io.ReadAll(req.Body)
		l.assertRequestBody(string(body))
	}
}

func (l testLogger) ResponseHttp(ctx context.Context, resp *http.Response, d time.Duration) {
	resp = sanitize.HTTPResponse(resp)
	if l.assertResponse != nil {
		l.assertResponse(resp)
	}
	if l.assertResponseBody != nil {
		body, _ := io.ReadAll(resp.Body)
		l.assertResponseBody(string(body))
		_ = resp.Body.Close()
	}
}
func (l testLogger) Error(ctx context.Context, err error) {}

func testOscClient(t *testing.T, l testLogger, resp any) *osc.Client {
	t.Helper()

	server := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	server.Start()
	cl, err := osc.NewClient(&profile.Profile{
		AccessKey: "foo",
		SecretKey: sk,
		Region:    "foo",
		Endpoints: profile.Endpoint{
			API: server.URL,
		},
	}, options.WithLogging(l))
	require.NoError(t, err)
	return cl
}

func testOksClient(t *testing.T, l testLogger, resp any) *oks.Client {
	t.Helper()

	server := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	server.Start()
	cl, err := oks.NewClient(&profile.Profile{
		AccessKey: "foo",
		SecretKey: sk,
		Region:    "foo",
		Endpoints: profile.Endpoint{
			OKS: server.URL,
		},
	}, options.WithLogging(l))
	require.NoError(t, err)
	return cl
}

func TestHTTPRequest(t *testing.T) {
	t.Run("pii fields are redacted", func(t *testing.T) {
		r := osc.CreateUserRequest{
			UserName:  "foo",
			UserEmail: new("john.doe@example.com"),
		}
		l := testLogger{
			assertRequestBody: func(body string) {
				assert.Equal(t, `{"UserEmail":"[REDACTED]","UserName":"[REDACTED]"}`, body) //nolint: testifylint
			},
		}
		cl := testOscClient(
			t, l,
			osc.CreateUserResponse{},
		)
		_, _ = cl.CreateUser(t.Context(), r)
	})
	t.Run("\\ does not cause issues", func(t *testing.T) {
		r := osc.CreateUserRequest{
			UserName:  "foo\"",
			UserEmail: new("john.doe@example\".com"),
		}
		l := testLogger{
			assertRequestBody: func(body string) {
				assert.Equal(t, `{"UserEmail":"[REDACTED]","UserName":"[REDACTED]"}`, body) //nolint: testifylint
			},
		}
		cl := testOscClient(
			t, l,
			osc.CreateUserResponse{},
		)
		_, _ = cl.CreateUser(t.Context(), r)
	})
}

func TestHTTPResponse(t *testing.T) {
	t.Run("osc fields are redacted", func(t *testing.T) {
		resp := osc.CreateAccessKeyResponse{
			AccessKey: &osc.AccessKeySecretKey{
				AccessKeyId: new("foo"),
				SecretKey:   new("bar"),
			},
		}
		l := testLogger{
			assertResponseBody: func(body string) {
				assert.Equal(t, `{"AccessKey":{"AccessKeyId":"foo","SecretKey":"[REDACTED]"}}`, strings.TrimSpace(body)) //nolint: testifylint
			},
		}
		cl := testOscClient(
			t, l,
			resp,
		)
		_, _ = cl.CreateAccessKey(t.Context(), osc.CreateAccessKeyRequest{})
	})
	t.Run("oks fields are redacted", func(t *testing.T) {
		resp := oks.KubeconfigResponse{
			Cluster: oks.ClustersClusterSchemaRPCResponse{
				Data: oks.KubeconfigData{
					Kubeconfig: "foo",
				},
			},
		}
		l := testLogger{
			assertRequest: func(req *http.Request) {
				assert.Equal(t, sanitize.Redacted, req.Header.Get("SecretKey"))
			},
			assertResponseBody: func(body string) {
				assert.Equal(t, `{"Cluster":{"data":{"kubeconfig":"[REDACTED]"},"request_id":""},"ResponseContext":{}}`, strings.TrimSpace(body)) //nolint: testifylint
			},
		}
		cl := testOksClient(
			t, l,
			resp,
		)
		_, _ = cl.GetKubeconfig(t.Context(), "foo", &oks.GetKubeconfigParams{})
	})
}
