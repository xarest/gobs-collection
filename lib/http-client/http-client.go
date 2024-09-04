package httpclient

import (
	"context"
	"io"
	"net/http"

	"github.com/xarest/gobs"
)

type Header map[string]string

type HTTPCient struct {
	c *http.Client
}

// Setup implements gobs.IServiceSetup.
func (h *HTTPCient) Setup(ctx context.Context, _ ...gobs.IService) error {
	h.c = http.DefaultClient
	return nil
}

func (h *HTTPCient) Get(ctx context.Context, uri string, header Header) ([]byte, error) {
	return h.Request(ctx, http.MethodGet, uri, nil, nil)
}

func (h *HTTPCient) Post(ctx context.Context, uri string, body io.Reader, header Header) ([]byte, error) {
	header["Content-Type"] = "application/json"
	return h.Request(ctx, http.MethodPost, uri, body, header)
}

func (h *HTTPCient) Request(ctx context.Context, method string, uri string, body io.Reader, header Header) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := h.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

var _ gobs.IServiceSetup = (*HTTPCient)(nil)
