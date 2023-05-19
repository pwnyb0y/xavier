package main

import (
	"context"
	"encoding/json"
	pb "github.com/pwnyb0y/xavier/gen/go/proto/xavier/v1"
	"io"
	"log"
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
		log.Printf("failed to make request to OpenAI API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return nil, err
	}

	log.Printf("response body: %s", body)

	var models []*pb.Model
	err = json.Unmarshal(body, &models)
	if err != nil {
		log.Printf("failed to unmarshal response body: %v", err)
		return nil, err
	}

	for _, model := range models {
		var permissions []*pb.Permission
		for _, perm := range model.Permissions {
			permissions = append(permissions, &pb.Permission{
				AllowCreateEngine:  perm.AllowCreateEngine,
				AllowSampling:      perm.AllowSampling,
				AllowLogprobs:      perm.AllowLogprobs,
				AllowSearchIndices: perm.AllowSearchIndices,
				AllowView:          perm.AllowView,
				AllowFineTuning:    perm.AllowFineTuning,
			})
		}
		models = append(models, &pb.Model{
			Id:          model.Id,
			Permissions: permissions,
		})
	}

	return &pb.GetModelsResponse{Models: models}, nil
}
