package main

import (
	"ai-test/gemini"
)

func main() {
	client := gemini.NewGeminiCient()
	gemini.RunPrompt(client, "Provide a working Go program that calls the Showcase API")
}
