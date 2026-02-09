package gemini

import (
	"ai-test/server/errors"
	"ai-test/server/responses"
	"ai-test/util"
	"ai-test/util/level"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
	Status responses.GenerationStatus
}

func NewClient() *Client {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  conf.Vertex.Project.Id,
		Location: conf.Vertex.Model.Location,
		Backend:  genai.BackendVertexAI,
	})

	util.HandleError("Failed to create client: %v", err, level.FATAL)

	configureAITools()

	return &Client{
		client: client,
	}
}

func (client *Client) RunCodeGenerationPrompt(prompt string) (*responses.GenerationResponse, *errors.HttpError) {
	ctx := context.Background()
	start := time.Now()

	client.Status = responses.Generating

	generatedResponse, err := client.client.Models.GenerateContent(
		ctx,
		conf.Vertex.Model.Name,
		genai.Text(prompt),
		groundedSearchConfig,
	)

	client.Status = responses.Generated

	log.Infof("Generated response in %v", time.Since(start))
	util.HandleError("Error generating response: %v", err, level.ERROR)

	var groundedText []string

	if len(generatedResponse.Candidates) == 0 || generatedResponse.Candidates[0].Content == nil {
		log.Error("No candidates found")
		return nil, &errors.InternalServerError
	}

	candidate := generatedResponse.Candidates[0]
	for _, part := range candidate.Content.Parts {
		groundedText = append(groundedText, part.Text)
	}

	return client.runJsonFormattingPrompt(strings.Join(groundedText, ""))
}

func (client *Client) runJsonFormattingPrompt(prompt string) (*responses.GenerationResponse, *errors.HttpError) {
	ctx := context.Background()
	start := time.Now()

	client.Status = responses.Formatting

	generatedResponse, err := client.client.Models.GenerateContent(
		ctx,
		conf.Vertex.Model.Name,
		genai.Text(fmt.Sprintf("Format this from markdown to json: %s", prompt)),
		formattingConfig,
	)

	log.Infof("Formatted response in %v", time.Since(start))
	util.HandleError("Error while formatting generated response: %v", err, level.ERROR)

	var formattedText []string

	if len(generatedResponse.Candidates) == 0 || generatedResponse.Candidates[0].Content == nil {
		log.Error("No candidates found")
		return nil, &errors.InternalServerError
	}

	candidate := generatedResponse.Candidates[0]
	for _, part := range candidate.Content.Parts {
		formattedText = append(formattedText, part.Text)
	}

	response := responses.GenerationResponse{}

	err = json.Unmarshal([]byte(strings.Join(formattedText, "")), &response)

	response.Time = time.Now()
	return &response, nil
}
