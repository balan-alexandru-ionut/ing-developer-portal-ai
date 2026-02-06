package main

import (
	"log"
	"net/http"
	"os"

	"ai-test/config"
	"ai-test/gemini"
)

func generateHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config.ReadConfigFile()

	// Create Gemini client
	client := gemini.NewGeminiClient()

	// Get generated JSON code
	result := gemini.RunPrompt(client, "Provide a working Java program that calls the Showcase API")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(result))
}

func main() {

	// ---- API route
	http.HandleFunc("/api/generate", generateHandler)

	// ---- Static Vue frontend (served from /dist after Vite build)
	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", fs)

	// ---- Cloud Run requires listening on $PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Local dev fallback
	}

	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
