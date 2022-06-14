package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	VERSION         = "version"
	DEFAULT_VERSION = "v1"
)

func getVersion() string {
	version, ok := os.LookupEnv(VERSION)
	if !ok {
		version = DEFAULT_VERSION
	}
	return version
}

func version(w http.ResponseWriter, req *http.Request) {
	//w.Header()Set("Content-Type", "application/json")
	fmt.Fprintln(w, getVersion())
}

func main() {
	http.HandleFunc("/version", version)
	http.ListenAndServe(":8090", nil)
}
