package routes

import "github.com/gofiber/fiber/v3"

func ConfigureRoutes(group *fiber.Router) {
	generateGroup := (*group).Group("generate")
	generateGroup.Get("/code", generateCode)
	generateGroup.Get("/status", generationStatus)
}
