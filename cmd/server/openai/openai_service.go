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

// OpenAIModelsResponse represents the structure of the response from the OpenAI API.
type OpenAIModelsResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Object     string `json:"object"`
		Created    int64  `json:"created"`
		OwnedBy    string `json:"owned_by"`
		Permission []struct {
			ID                 string `json:"id"`
			Object             string `json:"object"`
			Created            int64  `json:"created"`
			AllowCreateEngine  bool   `json:"allow_create_engine"`
			AllowSampling      bool   `json:"allow_sampling"`
			AllowLogprobs      bool   `json:"allow_logprobs"`
			AllowSearchIndices bool   `json:"allow_search_indices"`
			AllowView          bool   `json:"allow_view"`
			AllowFineTuning    bool   `json:"allow_fine_tuning"`
			Organization       string `json:"organization"`
			Group              string `json:"group"`
			IsBlocking         bool   `json:"is_blocking"`
		} `json:"permission"`
		Root   string `json:"root"`
		Parent string `json:"parent"`
	} `json:"data"`
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

// OpenAICompletionResponse represents the structure of the response from the OpenAI completion API.
type OpenAICompletionResponse struct {
	ID      string                `json:"id"`
	Object  string                `json:"object"`
	Created int64                 `json:"created"`
	Choices []pb.CompletionChoice `json:"choices"`
	Usage   pb.CompletionUsage    `json:"usage"`
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

	var completionResponse OpenAICompletionResponse
	err = json.Unmarshal(respBody, &completionResponse)
	if err != nil {
		log.Printf("failed to unmarshal response body: %v", err)
		return nil, err
	}

	var choices []*pb.CompletionChoice
	for _, choice := range completionResponse.Choices {
		choices = append(choices, &pb.CompletionChoice{
			Text:         choice.Text,
			Index:        choice.Index,
			Logprobs:     choice.Logprobs,
			FinishReason: choice.FinishReason,
		})
	}

	usage := &pb.CompletionUsage{
		PromptTokens:     completionResponse.Usage.PromptTokens,
		CompletionTokens: completionResponse.Usage.CompletionTokens,
		TotalTokens:      completionResponse.Usage.TotalTokens,
	}

	return &pb.CompletionResponse{
		Id:      completionResponse.ID,
		Object:  completionResponse.Object,
		Created: completionResponse.Created,
		Choices: choices,
		Usage:   usage,
	}, nil
}
