package routes

import (
	"ai-test/gemini"
	"ai-test/server/errors"
	"ai-test/server/responses"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"google.golang.org/genai"
)

var geminiClient *genai.Client

type PromptQuery struct {
	Prompt string `query:"prompt"`
}

func generateCode(c fiber.Ctx) {
	if geminiClient == nil {
		geminiClient = gemini.NewGeminiClient()
	}

	q := new(PromptQuery)

	if err := c.Bind().Query(q); err != nil {
		errors.BadRequestError.Send(c)
	}

	log.Info(q.Prompt)

	generatedCode, httpError := gemini.RunPrompt(geminiClient, q.Prompt)

	if httpError != nil {
		httpError.Send(c)
		return
	}

	if err := c.Status(http.StatusOK).JSON(generatedCode); err != nil {
		errors.InternalServerError.Send(c)
	}

	gemini.GenerationStatus = responses.Done
}

func generationStatus(c fiber.Ctx) {
	status := gemini.GenerationStatus

	if err := c.Status(http.StatusOK).JSON(responses.NewGenerationStatusResponse(status)); err != nil {
		errors.InternalServerError.Send(c)
	}
}
