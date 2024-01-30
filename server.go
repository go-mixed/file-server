package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

type mixedServer struct {
	dir         string
	fileServer  http.Handler
	proxyServer http.Handler
}

func newMixedServer(dir string, proxyModel bool) *mixedServer {
	s := &mixedServer{
		dir:        dir,
		fileServer: http.FileServer(http.Dir(dir)),
	}
	if proxyModel {
		s.proxyServer = newTransparentProxy()
	}
	return s
}

func (s *mixedServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.proxyServer != nil {
		var host string
		// 拥有X-Forwarded-Host头，直接转发
		if host = r.Header.Get("X-Forwarded-Host"); host != "" {
			s.proxyServer.ServeHTTP(w, r)
			return
		}

		// 从url path中分离host、path
		host, _ = parseUrlPath(r.URL.Path)
		if host != "" {
			// 文件、文件夹不存在，且域名合法
			if _, err := os.Stat(filepath.Join(s.dir, r.URL.Path)); err != nil {
				s.proxyServer.ServeHTTP(w, r)
				return
			}
		}
	}

	s.fileServer.ServeHTTP(w, r)
	return

}

// "/domain.com/path/to/file" -> "domain.com", "/path/to/file"
func parseUrlPath(originalUrlPath string) (host string, urlPath string) {
	urlPath = originalUrlPath
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}

	// 此处使用path.Clean而不是filepath.Clean，因为path.Clean不会将\替换为/
	filePath := path.Clean(urlPath)

	if segments := strings.Split(filePath, "/"); len(segments) > 1 {
		host = segments[1]
	}

	if host != "" && checkHost(host) == nil {
		return host, urlPath[len(host)+1:]
	}

	return "", urlPath
}

// Please use the package https://github.com/chmike/domain as is it maintained up to date with tests.

// checkHost returns an error if the domain name is not valid.
// See https://tools.ietf.org/html/rfc1034#section-3.5 and
// https://tools.ietf.org/html/rfc1123#section-2.
func checkHost(name string) error {
	// check for port
	if segments := strings.SplitN(name, ":", 2); len(segments) == 2 {
		name = segments[0]
		port, err := strconv.Atoi(segments[1])
		if err != nil {
			return fmt.Errorf("host has invalid port '%s'", segments[1])
		} else if port < 0 || port > 65535 {
			return fmt.Errorf("host has invalid port '%s', port must be between 0 and 65535", segments[1])
		}
	}

	switch {
	case len(name) == 0:
		return nil // an empty domain name will result in a cookie without a domain restriction
	case len(name) > 255:
		return fmt.Errorf("host name length is %d, can't exceed 255", len(name))
	}

	var l int
	for i := 0; i < len(name); i++ {
		b := name[i]
		if b == '.' {
			// check domain labels validity
			switch {
			case i == l:
				return fmt.Errorf("host has invalid character '.' at offset %d, label can't begin with a period", i)
			case i-l > 63:
				return fmt.Errorf("host byte length of label '%s' is %d, can't exceed 63", name[l:i], i-l)
			case name[l] == '-':
				return fmt.Errorf("host label '%s' at offset %d begins with a hyphen", name[l:i], l)
			case name[i-1] == '-':
				return fmt.Errorf("host label '%s' at offset %d ends with a hyphen", name[l:i], l)
			}
			l = i + 1
			continue
		}
		// test label character validity, note: tests are ordered by decreasing validity frequency
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-' || b >= 'A' && b <= 'Z') {
			// show the printable unicode character starting at byte offset i
			c, _ := utf8.DecodeRuneInString(name[i:])
			if c == utf8.RuneError {
				return fmt.Errorf("host has invalid rune at offset %d", i)
			}
			return fmt.Errorf("host has invalid character '%c' at offset %d", c, i)
		}
	}
	// check top level domain validity
	switch {
	case l == len(name):
		return fmt.Errorf("host has missing top level domain, domain can't end with a period")
	case len(name)-l > 63:
		return fmt.Errorf("host's top level domain '%s' has byte length %d, can't exceed 63", name[l:], len(name)-l)
	case name[l] == '-':
		return fmt.Errorf("host's top level domain '%s' at offset %d begin with a hyphen", name[l:], l)
	case name[len(name)-1] == '-':
		return fmt.Errorf("host's top level domain '%s' at offset %d ends with a hyphen", name[l:], l)
	case name[l] >= '0' && name[l] <= '9':
		return fmt.Errorf("host's top level domain '%s' at offset %d begins with a digit", name[l:], l)
	}
	return nil
}
