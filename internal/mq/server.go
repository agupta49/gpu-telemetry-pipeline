package mq

import (
	"context"
	"sync"

	"github.com/agupta49/gpu-telemetry-pipeline/pkg/pb"
)

type Server struct {
	pb.UnimplementedTelemetryQueueServer
	mu   sync.RWMutex
	subs []chan *pb.PublishRequest
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ch := range s.subs {
		select {
		case ch <- req:
		default: // drop if subscriber slow
		}
	}
	return &pb.PublishResponse{Success: true}, nil
}

func (s *Server) Subscribe(req *pb.PublishRequest, stream pb.TelemetryQueue_SubscribeServer) error {
	ch := make(chan *pb.PublishRequest, 100)
	s.mu.Lock()
	s.subs = append(s.subs, ch)
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		for i, c := range s.subs {
			if c == ch {
				s.subs = append(s.subs[:i], s.subs[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
		close(ch)
	}()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg := <-ch:
			if err := stream.Send(msg); err != nil {
				return err
			}
		}
	}
}
