package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/iter8-tools/iter8/abn/util"
	"github.com/iter8-tools/iter8/base/log"
	iter8abn "github.com/kalantar/ab-example/go/frontend/iter8"
)

const (
	// Default endpoint for ABn service
	ABN_SERVICE_ENDPOINT = "abn:50051"

	// Name of a backend service
	BACKEND_SERVICE = "backend"
	// Namespace where bacekend service deployed
	BACKEND_NAMESPACE = "default"
	// Default track to assign if ABn service can not find one
	DEFAULT_TRACK = "current"

	// Name of a sample metric
	SAMPLE_METRIC = "sample_metric"
)

var (
	// ABn service; assigned by init() method
	abnService iter8abn.DefaultABnService
)

// Assumes user is specified in a header X-User.
// If not set, a random user name will be assigned.
func user(req *http.Request) string {
	users, ok := req.Header["X-User"]
	if !ok {
		users = []string{util.RandomString(16)}
	}
	return users[0]
}

// getBackendURL returns a URL to the backend service to be used to satisffy a request for user
func getBackendURL(user string) string {
	// return "http://" + BACKEND_SERVICE + "-" + getTrack(user) + ":8090"
	return "http://" +
		BACKEND_SERVICE +
		"-" +
		abnService.GetTrack(user) +
		":8090"
}

// GET version of backend service
func version(w http.ResponseWriter, req *http.Request) {
	user := user(req)
	backendURL := getBackendURL(user)
	log.Logger.Infof("for user '%s', backendURL = %s", user, backendURL)

	resp, err := http.Get(backendURL + "/version")
	if err != nil {
		log.Logger.Info("GET /version failed: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Info("GET /version no data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(body))

	abnService.WriteMetric(SAMPLE_METRIC, strconv.Itoa(rand.Intn(100)), user)
}

func main() {
	http.HandleFunc("/version", version)
	http.ListenAndServe(":8091", nil)
}

func init() {
	abnService = iter8abn.DefaultABnService{
		AppName:      BACKEND_NAMESPACE + "/" + BACKEND_SERVICE,
		DefaultTrack: DEFAULT_TRACK,
		Service:      iter8abn.NewClient(ABN_SERVICE_ENDPOINT),
	}

}
