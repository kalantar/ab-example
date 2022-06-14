package iter8

import (
	"context"
	"time"

	pb "github.com/iter8-tools/iter8/abn/grpc"
	"github.com/iter8-tools/iter8/base/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ABnService interface {
	// Get recommended track for user
	GetTrack(user string) string
	// Write a metric value
	WriteMetric(name string, value string, user string)
}

// Default implementation of the ABnService interface
type DefaultABnService struct {
	// Name of backend application
	AppName string
	// Track to use if none can be identified (for example, the service is down)
	DefaultTrack string
	// gRPC client for ABn service
	Service *pb.ABNClient
}

// NewClient establishes a connection with the ABn service
func NewClient(endpoint string) *pb.ABNClient {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		log.Logger.Error("not able to connect to ABn service: ", err)
		return nil
	}

	client := pb.NewABNClient(conn)
	return &client
}

// Get recommended track for user
func (abn DefaultABnService) GetTrack(user string) string {
	track := abn.DefaultTrack
	if abn.Service == nil {
		log.Logger.Warn("ABn service not available -- track is ", track)
		return track
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s, err := (*abn.Service).Lookup(
		ctx,
		&pb.Application{
			Name: abn.AppName,
			User: user,
		},
	)
	if err != nil {
		log.Logger.Error("error -- track is ", track, " -- ", err)
		return track
	}

	track = s.GetTrack()
	log.Logger.Info("track is ", track)
	return track
}

// Write metric value
func (abn DefaultABnService) WriteMetric(name string, value string, user string) {
	if abn.Service == nil {
		log.Logger.Warn("ABn service not available")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := (*abn.Service).WriteMetric(
		ctx,
		&pb.MetricValue{
			Name:        name,
			Value:       value,
			Application: abn.AppName,
			User:        user,
		},
	)
	if err != nil {
		log.Logger.Error("error writing metric: ", err)
	}
}
