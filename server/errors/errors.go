package errors

import (
	"ai-test/server/responses"
	"ai-test/util"
	"ai-test/util/level"
	"time"

	"github.com/gofiber/fiber/v3"
)

var (
	BadRequestError = HttpError{
		Code:    400,
		Message: "There was an error with your request. Please try again.",
	}

	InternalServerError = HttpError{
		Code:    500,
		Message: "There was an error processing your request. Please try again later.",
	}
)

type HttpError struct {
	responses.HttpResponse
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *HttpError) Send(c fiber.Ctx) {
	e.Time = time.Now()
	if err := c.Status(e.Code).JSON(e); err != nil {
		util.HandleError("Could not send error response: %v", err, level.ERROR)
	}
}
