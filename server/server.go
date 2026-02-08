package server

import (
	"ai-test/server/routes"
	"ai-test/util"
	"ai-test/util/level"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

var app = fiber.New(fiber.Config{
	ServerHeader: "Fiber 3.0",
	AppName:      "ING AI Sandbox Playground",
})

func StartServer(port int) {
	api := app.Group("/api")
	routes.ConfigureRoutes(&api)

	app.Use("/", static.New("", static.Config{
		FS: os.DirFS("./frontend/dist"),
	}))

	app.Get("/*", func(c fiber.Ctx) error {
		return c.SendFile("./frontend/dist/index.html")
	})

	err := app.Listen(":"+strconv.Itoa(port), fiber.ListenConfig{
		EnablePrefork:     true,
		EnablePrintRoutes: true,
	})

	util.HandleError("Couldn't start server", err, level.FATAL)
}
