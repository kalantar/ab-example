package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	MY_VERSION      = "MY_VERSION"
	DEFAULT_VERSION = "v1"
)

// func getVersion() string {
// 	version, ok := os.LookupEnv(MY_VERSION)
// 	if !ok {
// 		version = DEFAULT_VERSION
// 	}
// 	return version
// }

type Data struct {
	Id   int
	Name string
}

// implment /recommend endpoint returning value of VERSION env variable
func recommend(w http.ResponseWriter, req *http.Request) {
	Logger.Trace("recommend called")
	// data := getVersion()
	data := Data{
		Id:   17,
		Name: "sample",
	}
	Logger.Info("/recommend returns ", data)
	Logger.Info(os.Environ())
	// fmt.Fprintln(w, data)
	json.NewEncoder(w).Encode(data)
}

var Logger *logrus.Logger

func main() {
	Logger = logrus.New()

	// configure backend service with "/recommend" endpoint
	http.HandleFunc("/recommend", recommend)
	http.ListenAndServe(":8091", nil)
}
