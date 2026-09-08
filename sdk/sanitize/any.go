package sanitize

import "net/http"

func Sanitize[T any](v T) T {
	vv := any(v)
	switch tv := vv.(type) {
	case string:
		vv = String(tv)
		return vv.(T)
	case *http.Request:
		vv = HTTPRequest(tv)
		return vv.(T)
	case *http.Response:
		vv = HTTPResponse(tv) //nolint: bodyclose
		return vv.(T)
	default:
		vv = Struct(vv)
		return vv.(T)
	}
}
