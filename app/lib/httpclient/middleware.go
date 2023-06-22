package httpclient

import (
	"net/http"
)

// MwBearerToken adds an Authorization header with the 'token'
// prefix.
//
// Example "Authorization: token abc123"
func MwBearerToken(token string) ClientMiddleware {
	return func(req *http.Request) (*http.Request, error) {
		req.Header.Add("Authorization", "token "+token)
		return req, nil
	}
}

func MwContentType(contentType string) ClientMiddleware {
	return func(req *http.Request) (*http.Request, error) {
		req.Header.Add("Content-Type", contentType)
		return req, nil
	}
}
