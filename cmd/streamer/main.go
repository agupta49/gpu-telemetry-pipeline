package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	streamer "github.com/agupta49/gpu-telemetry-pipeline/internal/streamer"
	"github.com/agupta49/gpu-telemetry-pipeline/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	csvPath := flag.String("csv", "data/dcgm_metrics_20250718_134233.csv", "CSV path")
	mqAddr := flag.String("mq", "mq:50051", "MQ address")
	flag.Parse()

	conn, err := grpc.Dial(*mqAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial mq: %v", err)
	}
	defer conn.Close()
	client := pb.NewMessageQueueClient(conn)

	for {
		if err := streamFile(*csvPath, client); err != nil {
			log.Printf("stream error: %v", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func streamFile(path string, client pb.MessageQueueClient) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ctx := context.Background()
	stream, err := client.Publish(ctx)
	if err != nil {
		return err
	}

	r := csv.NewReader(bufio.NewReader(file))
	if _, err := r.Read(); err != nil {
		return err
	}
	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		if !streamer.ValidateRecord(rec) {
			continue
		}
		tp := streamer.ParseCSV(rec)
		b, _ := json.Marshal(tp)
		if err := stream.Send(&pb.Message{Data: b, TimestampUnixNano: time.Now().UnixNano()}); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	_, err = stream.CloseAndRecv()
	return err
}
