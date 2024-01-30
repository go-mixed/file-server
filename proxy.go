package main

import (
	"net/http"
	"net/http/httputil"
	"strings"
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
	var host string
	var urlPath string

	// 拥有X-Forwarded-Host头，读取host、原URL path
	if host = r.Header.Get("X-Forwarded-Host"); host != "" {
		urlPath = r.URL.Path
	} else { // 从url path中分离host、path
		host, urlPath = parseUrlPath(r.URL.Path)
	}

	// host不包含端口，则从X-Forwarded-Port头中读取端口
	if !strings.Contains(host, ":") {
		if port := r.Header.Get("X-Forwarded-Port"); port != "" {
			host += ":" + port
		}
	}

	r.Host = host
	r.URL.Host = host
	r.URL.Path = urlPath

	// 通过X-Forwarded-Proto、X-Forwarded-Scheme设置scheme
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		r.URL.Scheme = proto
	} else if proto = r.Header.Get("X-Forwarded-Scheme"); proto != "" {
		r.URL.Scheme = proto
	} else {
		r.URL.Scheme = "https"
	}
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
