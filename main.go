package main

import (
	"flag"
	"net/http"
	"os"
)

func main() {
	var addr string
	var dir string
	flag.StringVar(&addr, "addr", ":8080", "listen address")
	flag.StringVar(&dir, "dir", ".", "directory to serve")
	flag.Parse()

	err := http.ListenAndServe(addr, CombinedLoggingHandler(os.Stdout, http.FileServer(http.Dir(dir))))
	if err != nil {
		panic(err)
	}
}
