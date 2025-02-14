package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy представляет собой обратный прокси-сервер
type ReverseProxy struct {
	proxy *httputil.ReverseProxy
}

// NewReverseProxy создаёт новый экземпляр ReverseProxy
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

// ServeHTTP обрабатывает запросы через прокси
func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}