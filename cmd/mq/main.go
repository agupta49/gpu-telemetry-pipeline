package main

import (
	"log"
	"net"

	"github.com/agupta49/gpu-telemetry-pipeline/internal/mq"
	"github.com/agupta49/gpu-telemetry-pipeline/pkg/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	svc := mq.NewServer()
	pb.RegisterTelemetryQueueServer(s, svc)
	log.Println("mq: starting on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
