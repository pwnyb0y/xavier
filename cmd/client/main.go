package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/zencodinglab/xavier/gen/go/proto/xavier/v1/openai"
	"google.golang.org/grpc"
)

const (
	serverAddress = "localhost:50051"
)

func main() {
	var prompt string
	flag.StringVar(&prompt, "prompt", "", "Prompt for completion")
	flag.Parse()

	if prompt == "" {
		log.Fatal("Prompt is required")
	}

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOpenAIClient(conn)

	// Prepare the request
	req := &pb.CompletionRequest{
		Model:       "davinci:ft-personal-2023-03-23-00-00-32",
		Prompt:      prompt,
		MaxTokens:   10,
		Temperature: 0.7,
	}

	// Call the Completion method
	resp, err := client.Completion(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Completion: %v", err)
	}

	// Print the response
	log.Printf("Response: %v", resp)
}
