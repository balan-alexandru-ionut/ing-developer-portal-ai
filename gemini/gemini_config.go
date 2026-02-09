package gemini

import (
	"ai-test/config"
	"fmt"

	"google.golang.org/genai"
)

var (
	conf                 = config.C
	groundedSearchConfig *genai.GenerateContentConfig
	formattingConfig     *genai.GenerateContentConfig
)

const systemPrompt = `You are a developer assistant.
                - Always use the provided documentation and API specifications to fulfill the user's requests.
                - All code should be generated for the Sandbox environment of ING. Always use the certificates, hosts and client ids for Sandbox.
                - The Sandbox host is always 'api.sandbox.ing.com'. For Sandbox access, no account creation on ING's developer portal is needed.
                - JWS signing is required only if the 'x-jws-signature' header is present on an endpoint.
                - Always use the endpoint names as they are, do not add versions or anything extra to them.
                - When an access token is required always refer to the OAuth 2.0 spec and documentation that were provided.
                - When an access token is required, the resulting code MUST contain the logic to obtain the access token.
                - All generated code should be functional so you need to always include mTLS setups or JWS/Cavage request signing setups. In the mTLS setup never verify the CA as it will never be provided.
                - Split the code into multiple files according to best practices and output a JSON where the keys are the file paths starting always from 'src/' and the values are the file contents.
                - Do not respond with anything else other than the JSON structure. Any text that guides the user on how to run the program should be put into a README.MD file under 'src/README.MD'.
                - Assume that certificates are always available under 'src/certs' directory.
                - If the user asks to change something only send back a JSON containing the files that need to be changed and the whole content of those files containing the requested changes.`

func configureAITools() {
	dataStorePath := fmt.Sprintf(
		"projects/%s/locations/%s/collections/default_collection/dataStores/%s",
		conf.Vertex.Project.Id,
		conf.Vertex.DataStore.Location,
		conf.Vertex.DataStore.Id,
	)

	searchTool := &genai.Tool{
		Retrieval: &genai.Retrieval{
			VertexAISearch: &genai.VertexAISearch{
				Datastore: dataStorePath,
			},
		},
	}

	systemInstructions := &genai.Content{
		Parts: []*genai.Part{
			{
				Text: systemPrompt,
			},
		},
	}

	groundedSearchConfig = &genai.GenerateContentConfig{
		SystemInstruction: systemInstructions,
		Tools:             []*genai.Tool{searchTool},
		Temperature:       &conf.Vertex.Model.Temperature,
		MaxOutputTokens:   conf.Vertex.Model.MaxOutputTokens,
		CandidateCount:    1,
	}

	formattingConfig = &genai.GenerateContentConfig{
		Temperature:      &conf.Vertex.Model.Temperature,
		MaxOutputTokens:  conf.Vertex.Model.MaxOutputTokens,
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"files": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"filePath": {Type: genai.TypeString},
							"code":     {Type: genai.TypeString},
						},
					},
					Required: []string{"filePath", "code"},
				},
			},
		},
	}
}
