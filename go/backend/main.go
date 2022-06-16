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

// implment /world endpoint returning value of VERSION env variable
func world(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, getVersion())
}

func main() {
	// configure backend service with "/world" endpoint
	http.HandleFunc("/world", world)
	http.ListenAndServe(":8090", nil)
}
