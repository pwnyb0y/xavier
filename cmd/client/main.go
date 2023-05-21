package main

import (
	"context"
	"flag"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	pb "github.com/zencodinglab/xavier/gen/go/proto/xavier/v1/openai"
	"google.golang.org/grpc"
)

const (
	serverAddress = "localhost:50051"
)

func main() {
	var prompt string
	flag.StringVar(&prompt, "prompt", "why did the chicken cross the AI?", "Prompt for completion")
	flag.Parse()

	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	client := pb.NewOpenAIClient(conn)

	// Prepare the request
	req := &pb.CompletionRequest{
		Model:       "text-davinci-003",
		Prompt:      prompt,
		MaxTokens:   100,
		Temperature: 0.7,
	}

	resp, err := client.Completion(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Completion: %v", err)
	}

	// Print the response
	log.Printf("Response: %v", resp)
}
