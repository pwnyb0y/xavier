package main

import (
	"context"
	"encoding/json"
	pb "github.com/pwnyb0y/xavier/gen/go/proto/xavier/v1"
	"io"
	"net/http"
)

// XavierServiceServer is a server for the Xavier service.
type XavierServiceServer struct {
	*pb.UnimplementedXavierServer
}

// OpenAIModelsResponse represents the structure of the response from the OpenAI API.
type OpenAIModelsResponse struct {
	Models []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"models"`
}

// GetModels reaches out to the OpenAI API and gets a list of all available models.
func (s *XavierServiceServer) GetModels(ctx context.Context, req *pb.GetModelsRequest) (*pb.GetModelsResponse, error) {
	config := LoadConfig()
	client := &http.Client{}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", "https://api.openai.com/v1/models", nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", "Bearer "+config.OpenAIKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var openAIResponse OpenAIModelsResponse
	err = json.Unmarshal(body, &openAIResponse)
	if err != nil {
		return nil, err
	}

	var models []*pb.Model
	for _, model := range openAIResponse.Models {
		models = append(models, &pb.Model{
			Id:   model.ID,
			Name: model.Name,
		})
	}

	return &pb.GetModelsResponse{Models: models}, nil
}
