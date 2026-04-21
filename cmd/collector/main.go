package main

import (
	"context"
	"flag"
	"log"

	"github.com/agupta49/gpu-telemetry-pipeline/internal/collector"
	"github.com/agupta49/gpu-telemetry-pipeline/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mqAddr := flag.String("mq", "localhost:50051", "MQ address")
	dbDSN := flag.String("db", "", "Postgres DSN")
	flag.Parse()

	log.Println("collector: starting")
	
	repo, err := collector.NewRepo(*dbDSN)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer repo.Close()

	conn, err := grpc.Dial(*mqAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("mq connect failed: %v", err)
	}
	defer conn.Close()
	
	client := pb.NewTelemetryQueueClient(conn)
	stream, err := client.Subscribe(context.Background(), &pb.PublishRequest{})
	if err != nil {
		log.Fatalf("subscribe failed: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("stream recv error: %v, reconnecting", err)
			return
		}
		if err := repo.Insert(msg.Point); err != nil {
			log.Printf("db insert error: %v", err)
		}
	}
}
