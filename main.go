package main

import (
	"flag"
	"net/http"
	"os"
)

func main() {
	var addr string
	var dir string
	var proxyModel bool
	flag.StringVar(&addr, "addr", ":8080", "listen address")
	flag.StringVar(&dir, "dir", ".", "directory to serve")
	flag.BoolVar(&proxyModel, "proxy", false, "proxy mode")
	flag.Parse()

	err := http.ListenAndServe(addr, CombinedLoggingHandler(os.Stdout, newMixedServer(dir, proxyModel)))
	if err != nil {
		panic(err)
	}
}
