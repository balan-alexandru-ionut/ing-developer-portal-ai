package main

import (
	"ai-test/config"
	"ai-test/gemini"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

)

func main() {
	config.ReadConfigFile()

	client := gemini.NewGeminiClient()
	gemini.RunPrompt(client, "Provide a working Java program that calls the Showcase API")
}
