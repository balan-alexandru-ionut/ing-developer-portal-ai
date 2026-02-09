package gemini

import (
	"ai-test/server/errors"
	"ai-test/server/responses"
	"ai-test/util"
	"ai-test/util/level"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
	Status responses.GenerationStatus
	chat   *genai.Chat
}

var Files []responses.GeneratedFile
var originalPrompt string

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
	originalPrompt = prompt

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

	groundedText := generatedResponse.Text()

	return client.runJsonFormattingPrompt(groundedText)
}

func (client *Client) StartChatSession() *errors.HttpError {
	ctx := context.Background()

	modelResponse, err := json.Marshal(Files)
	if err != nil {
		util.HandleError("Error marshalling model response: %v", err, level.ERROR)
		return &errors.InternalServerError
	}

	history := []*genai.Content{
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{Text: originalPrompt},
			},
		},
		{
			Role: genai.RoleModel,
			Parts: []*genai.Part{
				{Text: string(modelResponse)},
			},
		},
	}

	chatSession, err := client.client.Chats.Create(ctx, conf.Vertex.Model.Name, nil, history)
	if err != nil {
		util.HandleError("Error creating chat session: %v", err, level.ERROR)
		return &errors.InternalServerError
	}

	client.chat = chatSession

	return nil
}

func (client *Client) SendMessage(message string) (*responses.ChatResponse, *errors.HttpError) {
	ctx := context.Background()

	prompt := genai.Part{
		Text: message,
	}

	response, err := client.chat.SendMessage(ctx, prompt)
	if err != nil {
		return nil, &errors.InternalServerError
	}

	return responses.NewChatResponse(response.Text()), nil
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

	formattedText := generatedResponse.Text()

	response := responses.GenerationResponse{}

	err = json.Unmarshal([]byte(formattedText), &response)

	Files = response.Files

	response.Time = time.Now()
	return &response, nil
}
