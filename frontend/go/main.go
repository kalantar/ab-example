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
	// map of track to route to backend service
	trackToRoute = map[string]string{
		"current":   "http://backend-current:8091",
		"candidate": "http://backend-candidate:8091",
	}
)

// implment /getRecommendation endpoint
// calls backend service /recommend endpoint
func getRecommendation(w http.ResponseWriter, req *http.Request) {
	// Get user (session) identifier, for example by inspection of header X-User
	user := req.Header["X-User"][0]

	// Get endpoint of backend endpoint "/recommend"
	// In this example, the backend endpoint depends on the version (track) of the backend service
	// the user is assigned by the Iter8 SDK Lookup() method

	// establish connection to ABn service
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial("abn:50051", opts...)
	if err != nil {
		http.Error(w, fmt.Sprintf("error connecting to ABn service: %s", err), http.StatusInternalServerError)
		return
	}
	client := pb.NewABNClient(conn)

	// call ABn service API Lookup() to get an assigned track for the user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s, err := client.Lookup(
		ctx,
		&pb.Application{
			Name: "default/backend",
			User: user,
		},
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("ABn service Lookup() failed %s", err), http.StatusInternalServerError)
		return
	}

	// lookup route using track
	route := trackToRoute[s.GetTrack()]

	// call backend service using url
	resp, err := http.Get(route + "/recommend")
	if err != nil {
		http.Error(w, "call to backend endpoint /recommend failed", http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("backend endpoint /recommend returned no data %s", err), http.StatusInternalServerError)
		return
	}

	// write response to query
	fmt.Fprintln(w, "Recommendation: "+string(body))
}

// implment /buy endpoint
// writes value for sample_metric which may have spanned several calls to /getRecommendatoon
func buy(w http.ResponseWriter, req *http.Request) {
	// Get user (session) identifier, for example by inspection of header X-User
	user := req.Header["X-User"][0]

	// export metric to metrics database
	// this is best effort; we ignore any failure

	// establish connection to ABn service
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial("abn:50051", opts...)
	if err != nil {
		http.Error(w, fmt.Sprintf("error connecting to ABn service: %s", err), http.StatusInternalServerError)
		return
	}
	client := pb.NewABNClient(conn)

	// export metric
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _ = client.WriteMetric(
		ctx,
		&pb.MetricValue{
			Name:        "sample_metric",
			Value:       strconv.Itoa(rand.Intn(100)),
			Application: "default/backend",
			User:        user,
		},
	)
}

func main() {
	// configure frontend service with "/hello" and "/goodbye" endpoints
	http.HandleFunc("/getRecommendation", getRecommendation)
	http.HandleFunc("/buy", buy)
	http.ListenAndServe(":8090", nil)
}
