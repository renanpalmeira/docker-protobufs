package v1

import (
	"context"
	pb "github.com/renanpalmeira/docker-protobufs/api/v1/genproto"
)

type health struct {
}

func NewHealth() pb.HealthServer {
	return &health{}
}

func (h health) Check(_ context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

func (h health) Watch(request *pb.HealthCheckRequest, server pb.Health_WatchServer) error {
	return nil
}
