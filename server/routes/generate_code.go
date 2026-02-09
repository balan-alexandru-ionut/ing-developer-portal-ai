package routes

import (
	"ai-test/gemini"
	"ai-test/server/errors"
	"ai-test/server/responses"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

var geminiClient *gemini.Client

type PromptQuery struct {
	Prompt string `query:"prompt"`
}

func generateCode(c fiber.Ctx) {
	if geminiClient == nil {
		geminiClient = gemini.NewClient()
	}

	q := new(PromptQuery)

	if err := c.Bind().Query(q); err != nil {
		errors.BadRequestError.Send(c)
	}

	log.Info(q.Prompt)

	generatedCode, httpError := geminiClient.RunCodeGenerationPrompt(q.Prompt)

	if httpError != nil {
		httpError.Send(c)
		return
	}

	if err := c.Status(http.StatusOK).JSON(generatedCode); err != nil {
		errors.InternalServerError.Send(c)
	}

	geminiClient.Status = responses.Done
}

func generationStatus(c fiber.Ctx) {
	if geminiClient == nil {
		c.Status(http.StatusOK).JSON(responses.NewGenerationStatusResponse(responses.NotStarted))
		return
	}

	if err := c.Status(http.StatusOK).JSON(responses.NewGenerationStatusResponse(geminiClient.Status)); err != nil {
		errors.InternalServerError.Send(c)
	}
}
