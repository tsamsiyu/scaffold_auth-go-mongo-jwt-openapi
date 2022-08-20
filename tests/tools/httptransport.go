package tools

import (
	"net/http"
)

type MyHttpTransport struct {
	Scheme string
	Host   string

	BaseTransport http.RoundTripper
}

func (t *MyHttpTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	scheme := t.Scheme
	if scheme == "" {
		scheme = "http"
	}

	r.URL.Scheme = scheme
	r.URL.Host = t.Host

	return t.BaseTransport.RoundTrip(r)
}
