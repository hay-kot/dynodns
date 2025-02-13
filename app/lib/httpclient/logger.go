package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Transport struct {
	Transport http.RoundTripper
}

func NewTransport() *Transport {
	return &Transport{Transport: http.DefaultTransport}
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.logRequest(req)

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		LogErrRequest(log.Logger, resp, err)
		return nil, err
	}

	t.logResponse(resp)
	return resp, nil
}

func (t Transport) logRequest(req *http.Request) {
	log.Info().Str("method", req.Method).Str("url", req.URL.String()).Msg("req  ->")
}

func (t Transport) logResponse(resp *http.Response) {
	l := log.Info()

	if resp.StatusCode >= 400 {
		l = log.Error()
	}

	l.Str("method", resp.Request.Method).
		Str("status", resp.Status).
		Str("url", resp.Request.URL.String()).
		Msg("resp <-")
}

func LogErrRequest(l zerolog.Logger, req *http.Response, err error) {
	bodyCopy := new(bytes.Buffer)
	if req == nil {
		l.Error().Err(err).Msg("request is nil")
		return
	}

	if req.Body != nil {
		_, _ = bodyCopy.ReadFrom(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(bodyCopy.Bytes()))
	}

	l.Error().Str("method", req.Request.Method).Str("url", req.Request.URL.String()).Msg("request")
	l.Error().Msg("------- Begin Response -------")
	fmt.Println(fmtJSON(bodyCopy.Bytes()))
	l.Error().Msg("-------- End Body --------")
	l.Error().Err(err).Msgf("%s %s", req.Request.Method, req.Request.URL)
}

func fmtJSON(body []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, body, "", "  ")
	if err != nil {
		return string(body)
	}
	return out.String()
}
