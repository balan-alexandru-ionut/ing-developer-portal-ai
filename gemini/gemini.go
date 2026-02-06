package gemini

import (
    "ai-test/config"
    "ai-test/logger"
    "ai-test/util"
    "ai-test/util/level"
    "context"
    "fmt"
    "strings"
    "time"

    "google.golang.org/genai"
)

var log = logger.NewLogger()
var conf = config.C

var groundedSearchConfig *genai.GenerateContentConfig
var formattingConfig *genai.GenerateContentConfig

func NewGeminiClient() *genai.Client {
    ctx := context.Background()

    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        Project:  conf.Vertex.Project.Id,
        Location: conf.Vertex.Model.Location,
        Backend:  genai.BackendVertexAI,
    })

    util.HandleError("Failed to create client: %v", err, level.FATAL)

    return client
}

func RunPrompt(client *genai.Client, prompt string) string {
    configureAITools()

    ctx := context.Background()
    start := time.Now()

    generatedResponse, err := client.Models.GenerateContent(
        ctx,
        conf.Vertex.Model.Name,
        genai.Text(prompt),
        groundedSearchConfig,
    )

    log.Infof("Generated response in %v", time.Since(start))
    util.HandleError("Error generating response: %v", err, level.ERROR)

    groundedText := []string{}

    if len(generatedResponse.Candidates) == 0 || generatedResponse.Candidates[0].Content == nil {
        log.Error("No candidates found")
        return ""
    }

    candidate := generatedResponse.Candidates[0]
    for _, part := range candidate.Content.Parts {
        groundedText = append(groundedText, part.Text)
    }

    return runJsonFormattingPrompt(client, strings.Join(groundedText, ""))
}

func runJsonFormattingPrompt(client *genai.Client, text string) string {
    ctx := context.Background()
    start := time.Now()

    generatedResponse, err := client.Models.GenerateContent(
        ctx,
        conf.Vertex.Model.Name,
        genai.Text(fmt.Sprintf("Format this from markdown to json: %s", text)),
        formattingConfig,
    )

    log.Infof("Formatted response in %v", time.Since(start))
    util.HandleError("Error while formatting generated response: %v", err, level.ERROR)

    formattedText := []string{}

    if len(generatedResponse.Candidates) == 0 || generatedResponse.Candidates[0].Content == nil {
        log.Error("No candidates found")
        return ""
    }

    candidate := generatedResponse.Candidates[0]
    for _, part := range candidate.Content.Parts {
        formattedText = append(formattedText, part.Text)
    }

    return strings.Join(formattedText, "")
}

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
                Text: `You are a developer assistant.
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
                - If the user asks to change something only send back a JSON containing the files that need to be changed and the whole content of those files containing the requested changes.`,
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
