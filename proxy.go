package main

import (
	"net/http"
	"net/http/httputil"
)

// TransparentProxy 透明代理模式
type TransparentProxy struct {
	proxy *httputil.ReverseProxy
}

func newTransparentProxy() *TransparentProxy {
	p := &TransparentProxy{}

	p.proxy = &httputil.ReverseProxy{
		Rewrite:        nil,
		Director:       p.modifyRequest,
		Transport:      nil,
		FlushInterval:  0,
		ErrorLog:       nil,
		BufferPool:     nil,
		ModifyResponse: p.modifyResponse,
		ErrorHandler:   nil,
	}
	return p
}

func (t *TransparentProxy) modifyRequest(r *http.Request) {
	domain, upath := parseUrlPath(r.URL.Path)
	r.Host = domain
	r.URL.Host = domain
	r.URL.Path = upath
	r.URL.Scheme = "https"
	return
}

func (t *TransparentProxy) modifyResponse(response *http.Response) error {
	// 添加跨域头
	if response.Header.Get("Access-Control-Allow-Origin") != "*" {
		response.Header.Set("Access-Control-Allow-Origin", "*")
	}
	if response.Header.Get("Access-Control-Allow-Methods") == "" {
		response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, HEAD, PATCH")
	}
	if response.Header.Get("Access-Control-Allow-Headers") == "" {
		response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	if response.Header.Get("Access-Control-Allow-Credentials") == "" {
		response.Header.Set("Access-Control-Allow-Credentials", "true")
	}
	return nil
}

func (t *TransparentProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.proxy.ServeHTTP(w, r)
}
