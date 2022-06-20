package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	VERSION         = "VERSION"
	DEFAULT_VERSION = "v1"
)

func getVersion() string {
	version, ok := os.LookupEnv(VERSION)
	if !ok {
		version = DEFAULT_VERSION
	}
	return version
}

// implment /recommend endpoint returning value of VERSION env variable
func recommend(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, getVersion())
}

func main() {
	// configure backend service with "/recommend" endpoint
	http.HandleFunc("/recommend", recommend)
	http.ListenAndServe(":8091", nil)
}
