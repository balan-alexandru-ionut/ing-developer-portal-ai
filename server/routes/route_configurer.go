package routes

import "github.com/gofiber/fiber/v3"

func ConfigureRoutes(group *fiber.Router) {
	generateGroup := (*group).Group("generate")
	chatGroup := (*group).Group("chat")

	generateGroup.Get("/code", generateCode)
	generateGroup.Get("/status", generationStatus)
	generateGroup.Get("/archive", generateFilesHandler)

	chatGroup.Post("/start", startChat)
	chatGroup.Post("/message", chat)
}
