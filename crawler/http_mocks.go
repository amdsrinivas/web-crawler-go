package crawler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

type MockRoundTripper struct {
}

func (m MockRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {

	if request.URL.Host == "amazon.in" {
		return nil, errors.New("timeout")
	}
	return &http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte(request.URL.Host))),
	}, nil
}
