package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	pb "github.com/kalantar/ab-example/frontend/go/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// ABn service; assigned by init() method
	abnService   *pb.ABNClient
	trackToRoute = map[string]string{
		"current":   "http://backend-current:8090",
		"candidate": "http://backend-candidate:8090",
	}
)

// implment /hello endpoint
// calls backend service /version endpoint
func hello(w http.ResponseWriter, req *http.Request) {
	// Get user (session) identifier, for example by inspection of header X-User
	users, ok := req.Header["X-User"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "header X-User missing")
		return
	}
	user := users[0]

	// Get endpoint of backend endpoint "/world"
	// In this example, the backend endpoint depends on the version (track) of the backend service
	// the user is assigned by the Iter8 SDK Lookup() method

	// verify the ABn service is avaiable
	if abnService == nil {
		http.Error(w, "ABn service unavailable", http.StatusInternalServerError)
		return
	}

	// call ABn service API Lookup() to get an assigned track for the user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	s, err := (*abnService).Lookup(
		ctx,
		&pb.Application{
			Name: "default/backend",
			User: user,
		},
	)
	cancel()
	if err != nil {
		http.Error(w, fmt.Sprintf("ABn service Lookup() failed %s", err), http.StatusInternalServerError)
		return
	}

	// lookup route using track
	route, ok := trackToRoute[s.GetTrack()]
	if !ok {
		http.Error(w, fmt.Sprintf("unknown track returned: %s", s.GetTrack()), http.StatusInternalServerError)
		return
	}

	// call backend service using url
	resp, err := http.Get(route + "/world")
	if err != nil {
		http.Error(w, "call to backend endpoint /world failed", http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("backend endpoint /world returned no data %s", err), http.StatusInternalServerError)
		return
	}

	// write response to query
	fmt.Fprintln(w, "Hello world "+string(body))

	// export metric to metrics database
	// this is best effort; we ignore any failure
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	_, _ = (*abnService).WriteMetric(
		ctx,
		&pb.MetricValue{
			Name:        "sample_metric",
			Value:       strconv.Itoa(rand.Intn(100)),
			Application: "default/backend",
			User:        user,
		},
	)
	cancel()
}

func main() {
	// establish connect to ABn service
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial("abn:50051", opts...)
	if err != nil {
		fmt.Printf("unable to connect to ABn service: %s\n", err.Error())
		return
	}

	client := pb.NewABNClient(conn)
	abnService = &client

	// configure frontend service with "/hello" endpoint
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8091", nil)
}