package routes

import (
	"ai-test/gemini"
	"ai-test/server/errors"
	"ai-test/server/responses"
	"ai-test/util"
	"ai-test/util/level"

	"github.com/gofiber/fiber/v3"
)

type ChatPrompt struct {
	Prompt string `json:"prompt"`
}

func startChat(c fiber.Ctx) {
	if geminiClient == nil {
		geminiClient = gemini.NewClient()
	}

	if httpErr := geminiClient.StartChatSession(); httpErr != nil {
		httpErr.Send(c)
		return
	}

	if err := c.Status(200).SendString("Chat session started"); err != nil {
		errors.InternalServerError.Send(c)
	}
}

func chat(c fiber.Ctx) {
	if geminiClient == nil {
		herr := &errors.HttpError{
			HttpResponse: responses.HttpResponse{}.Zero(),
			Code:         400,
			Message:      "Chat must be initialized first",
		}
		herr.Send(c)
		return
	}

	chatPrompt := new(ChatPrompt)
	if err := c.Bind().JSON(chatPrompt); err != nil {
		util.HandleError(err.Error(), err, level.WARN)
		errors.BadRequestError.Send(c)
		return
	}

	response, herr := geminiClient.SendMessage(chatPrompt.Prompt)
	if herr != nil {
		herr.Send(c)
		return
	}

	if err := c.Status(200).JSON(response); err != nil {
		errors.InternalServerError.Send(c)
	}
}
