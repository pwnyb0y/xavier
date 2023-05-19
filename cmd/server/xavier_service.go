package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	pb "github.com/pwnyb0y/xavier/gen/go/proto/xavier/v1"
)

// XavierServiceServer is a server for the Xavier service.
type XavierServiceServer struct {
	*pb.UnimplementedXavierServer
}

// OpenAIModel represents the structure of a model in the OpenAI API response.
type OpenAIModel struct {
	ID          string             `json:"id"`
	Object      string             `json:"object"`
	Created     int64              `json:"created"`
	OwnedBy     string             `json:"owned_by"`
	Permissions []OpenAIPermission `json:"permissions"`
	Root        string             `json:"root"`
	Parent      string             `json:"parent"`
}

// OpenAIPermission represents the structure of a permission in the OpenAI API response.
type OpenAIPermission struct {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return nil, err
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("failed to close response body: %v", closeErr)
		}
	}()

	var openAIResponse struct {
		Models []OpenAIModel `json:"models"`
	}
	err = json.Unmarshal(body, &openAIResponse)
	if err != nil {
		log.Printf("failed to unmarshal response body: %v", err)
		return nil, err
	}

	var models []*pb.Model
	for _, model := range openAIResponse.Models {
		var permissions []*pb.Permission
		for _, perm := range model.Permissions {
			permissions = append(permissions, &pb.Permission{
				Id:                 perm.ID,
				Object:             perm.Object,
				Created:            perm.Created,
				AllowCreateEngine:  perm.AllowCreateEngine,
				AllowSampling:      perm.AllowSampling,
				AllowLogprobs:      perm.AllowLogprobs,
				AllowSearchIndices: perm.AllowSearchIndices,
				AllowView:          perm.AllowView,
				AllowFineTuning:    perm.AllowFineTuning,
				Organization:       perm.Organization,
				Group:              perm.Group,
				IsBlocking:         perm.IsBlocking,
			})
		}
		models = append(models, &pb.Model{
			Id:          model.ID,
			Permissions: permissions,
			Root:        model.Root,
			Parent:      model.Parent,
		})
	}

	return &pb.GetModelsResponse{Models: models}, nil
}
