package main

import (
	"log"
	"net"

	pb "github.com/pwnyb0y/xavier/gen/go/proto/xavier/v1"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	server := &XavierServiceServer{}
	pb.RegisterXavierServer(s, server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
