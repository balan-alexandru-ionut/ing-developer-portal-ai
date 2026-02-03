package main

import (
	"ai-test/config"
	"ai-test/gemini"
)

func main() {
	config.ReadConfigFile()

	client := gemini.NewGeminiCient()
	gemini.RunPrompt(client, "Provide a working Go program that calls the Showcase API")
}
