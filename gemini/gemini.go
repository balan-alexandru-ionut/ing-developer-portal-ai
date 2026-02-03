package gemini

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

const (
	project_id          = "ing-dev-portal"
	location_europe     = "europe-central2"
	location_data_store = "global"
	data_store_id       = "dev-portal-store_1770036213809"
)

var temperature = float32(0.2)
var max_output_tokens = int32(16384)

var dataStorePath = fmt.Sprintf("projects/%s/locations/%s/collections/default_collection/dataStores/%s", project_id, location_data_store, data_store_id)
var searchTool = &genai.Tool{
	Retrieval: &genai.Retrieval{
		VertexAISearch: &genai.VertexAISearch{
			Datastore: dataStorePath,
		},
	},
}
var systemInstructions = &genai.Content{
	Parts: []*genai.Part{
		{
			Text: `You are a developer assistant.
				- Always use the provided documentation and API specifications to fulfill the user's requests.
				- All code should be generated for the Sandbox environment of ING. Always use the certificates, hosts and client ids for Sandbox.
				- The Sandbox host is always 'api.sandbox.ing.com'. For Sandbox access, no account creation on ING's developer portal is needed.
				- JWS signing is required only if the 'x-jws-signature' header is present on an endpoint.
				- Always use the endpoint names as they are, do not add versions or anything extra to them.
				- When an access token is required always refer to the OAuth 2.0 spec and documentation that were provided.
				- When an access token is required, the resulting code MUST contain the logic to obtain the access token.
				- All generated code should be functional so you need to always include mTLS setups or JWS/Cavage request signing setups.
				- Split the code into multiple files according to best practices and output a JSON where the keys are the filepaths starting always from 'src/' and the values are the file contents.
				- Do not respond with anything else other than the JSON structure. Any text that guides the user on how to run the program should be put into a README.MD file under 'src/README.MD'.
				- Assume that certificates are always available under 'src/certs' directory.
				- If the user asks to change something only send back a JSON containing the files that need to be changed and the whole content of those files containing the requested changes.`,
		},
	},
}
var generateContentConfig = &genai.GenerateContentConfig{
	SystemInstruction: systemInstructions,
	Tools:             []*genai.Tool{searchTool},
	Temperature:       &temperature,
	MaxOutputTokens:   max_output_tokens,
}

func NewGeminiCient() *genai.Client {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  project_id,
		Location: location_europe,
		Backend:  genai.BackendVertexAI,
	})

	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	return client
}

func RunPrompt(client *genai.Client, prompt string) {
	ctx := context.Background()

	responseIterator := client.Models.GenerateContentStream(ctx, "gemini-2.5-flash", genai.Text(prompt), generateContentConfig)

	for chunk, err := range responseIterator {
		if err != nil {
			log.Fatalf("Error while streaming response: %v\n", err)
		}

		for _, candidate := range chunk.Candidates {
			if candidate.Content != nil {
				for _, part := range candidate.Content.Parts {
					fmt.Print(part.Text)
				}
			}
		}
	}
}
