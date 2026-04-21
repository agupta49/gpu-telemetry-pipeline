package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/agupta49/gpu-telemetry-pipeline/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	csvPath := flag.String("csv", "", "path to DCGM csv")
	mqAddr := flag.String("mq", "localhost:50051", "MQ address")
	flag.Parse()

	log.Println("streamer: starting")
	
	conn, err := grpc.Dial(*mqAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewTelemetryQueueClient(conn)

	for {
		if _, err := os.Stat(*csvPath); os.IsNotExist(err) {
			log.Printf("streamer: %s not found, retrying in 5s", *csvPath)
			time.Sleep(5 * time.Second)
			continue
		}
		streamFile(*csvPath, client)
		log.Println("streamer: finished CSV, looping")
		time.Sleep(1 * time.Second)
	}
}

func streamFile(path string, client pb.TelemetryQueueClient) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("error opening csv: %v", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum == 1 { continue } // skip header
		fields := strings.Split(scanner.Text(), ",")
		if len(fields) < 3 { continue }
		
		gpuID := fields[0]
		metricName := fields[1]
		val, _ := strconv.ParseFloat(fields[2], 64)
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := client.Publish(ctx, &pb.PublishRequest{
			Point: &pb.TelemetryPoint{
				GpuId:      gpuID,
				MetricName: metricName,
				Value:      val,
				Timestamp:  time.Now().Unix(), // Per requirement: use processing time
			},
		})
		cancel()
		if err != nil {
			log.Printf("publish error: %v", err)
			time.Sleep(1 * time.Second)
		}
		time.Sleep(100 * time.Millisecond) // simulate stream rate
	}
}
