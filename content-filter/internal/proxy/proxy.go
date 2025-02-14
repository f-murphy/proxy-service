package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	proxy *httputil.ReverseProxy
}

func NewReverseProxy(target string) *ReverseProxy {
	targetURL, _ := url.Parse(target)
	return &ReverseProxy{
		proxy: &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = targetURL.Scheme
				req.URL.Host = targetURL.Host
				req.Host = targetURL.Host
			},
		},
	}
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}