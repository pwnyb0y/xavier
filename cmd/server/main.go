package main

import (
	"log"
	"net"

	pb "github.com/pwnyb0y/xavier/gen/go/proto/xavier/v1"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting Xavier server...")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Successfully started listening on localhost:50051")

	s := grpc.NewServer()
	server := &XavierServiceServer{}
	pb.RegisterXavierServer(s, server)

	log.Println("Xavier server registered, ready to serve requests")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Println("Xavier server stopped")
}
