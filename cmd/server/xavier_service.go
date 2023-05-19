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
func (s *XavierServiceServer) GetModels(ctx context.Context, req *pb.GetModelsRequest) (*pb.GetModelsResponse, error) {
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

	log.Printf("Response body: %v", string(body))
	var openAIResponse OpenAIModelsResponse
	err = json.Unmarshal(body, &openAIResponse)
	if err != nil {
		log.Printf("failed to unmarshal response body: %v", err)
		return nil, err
	}

	var models []*pb.Model
	for _, model := range openAIResponse.Data {
		var permissions []*pb.Permission
		for _, perm := range model.Permission {
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
			Object:      model.Object,
			Created:     model.Created,
			OwnedBy:     model.OwnedBy,
			Permissions: permissions,
			Root:        model.Root,
			Parent:      model.Parent,
		})
	}

	return &pb.GetModelsResponse{Models: models}, nil
}
