package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	pb "github.com/zencodinglab/xavier/gen/go/proto/xavier/v1/openai"
)

// OpenAIServiceServer is a server for the OpenAI service.
type OpenAIServiceServer struct {
	*pb.UnimplementedOpenAIServer
}

// GetModels reaches out to the OpenAI API and gets a list of all available models.
func (s *OpenAIServiceServer) GetModels(ctx context.Context, req *pb.GetModelsRequest) (*pb.GetModelsResponse, error) {
	log.Printf("Received request: %v", req)
	config := LoadConfig()
	client := &http.Client{}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", "https://api.openai.com/v1/models", nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", "Bearer "+config.OpenAIKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("failed to make request to OpenAI API: %v", err)
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return nil, err
	}

	var openaiResponse *pb.OpenAIGetModelsResponse
	err = json.Unmarshal(body, &openaiResponse)
	if err != nil {
		log.Printf("failed to unmarshal response body: %v", err)
		return nil, err
	}

	response := &pb.GetModelsResponse{Models: openaiResponse.Data}
	return response, nil
}

// Completion reaches out to the OpenAI API and performs a GPT-4 completion.
func (s *OpenAIServiceServer) Completion(ctx context.Context, req *pb.CompletionRequest) (*pb.CompletionResponse, error) {
	log.Printf("Received completion request: %v", req)
	config := LoadConfig()
	client := &http.Client{}
	body, err := json.Marshal(req)
	if err != nil {
		log.Printf("failed to marshal request body: %v", err)
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/completions", bytes.NewReader(body))
	if err != nil {
		log.Printf("failed to create HTTP request: %v", err)
		return nil, err
	}

	httpReq.Header.Add("Authorization", "Bearer "+config.OpenAIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("failed to make request to OpenAI API: %v", err)
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return nil, err
	}

	log.Printf("respBody: %v", string(respBody))
	var completionResponse *pb.CompletionResponse
	err = json.Unmarshal(respBody, &completionResponse)
	if err != nil {
		log.Printf("failed to unmarshal response body: %v", err)
		return nil, err
	}

	return completionResponse, nil
}
