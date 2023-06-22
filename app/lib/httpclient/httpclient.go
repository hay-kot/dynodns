// Package httpclient provides a client for working with HTTP requests, using
// a method based API.
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ClientMiddleware = func(*http.Request) (*http.Request, error)

type Client struct {
	client  *http.Client
	baseURL string
	mw      []ClientMiddleware
}

func New(client *http.Client, base string) *Client {
	return &Client{
		client:  client,
		baseURL: base,
		mw:      nil,
	}
}

func (c *Client) Use(mws ...ClientMiddleware) {
	c.mw = append(c.mw, mws...)
}

func (c *Client) Get(url string, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) GetCtx(ctx context.Context, url string, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Post(url string, payload []byte, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) PostCtx(ctx context.Context, url string, payload []byte, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Put(url string, payload []byte, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) PutCtx(ctx context.Context, url string, payload []byte, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Delete(url string, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) DeleteCtx(ctx context.Context, url string, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Do(req *http.Request, rmw []ClientMiddleware) (*http.Response, error) {
	for _, mw := range c.mw {
		var err error
		req, err = mw(req)
		if err != nil {
			return nil, err
		}
	}

	for _, mw := range rmw {
		var err error
		req, err = mw(req)
		if err != nil {
			return nil, err
		}
	}

	return c.client.Do(req)
}

// Path will safely join the base URL and the provided path and return a string
// that can be used in a request.
func (c *Client) Path(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}

	base := strings.TrimRight(c.baseURL, "/")
	if url == "" {
		return base
	}

	return base + "/" + strings.TrimLeft(url, "/")
}

// Pathf will call fmt.Sprintf with the provided values and then pass them
// to Client.Path as a convenience.
func (c *Client) Pathf(url string, v ...any) string {
	url = fmt.Sprintf(url, v...)
	return c.Path(url)
}

func Decode(r *http.Response, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(val); err != nil {
		return err
	}
	return nil
}
