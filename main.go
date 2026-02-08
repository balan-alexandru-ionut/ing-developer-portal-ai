package main

import (
	"ai-test/config"
	"ai-test/server"
)

func main() {
	config.ReadConfigFile()
	server.StartServer(8080)
}
