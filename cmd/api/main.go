package main

import (
	"log"
	"os"
	"os/signal"
	"qualifire-home-assignment/internal/configs"
	"qualifire-home-assignment/internal/http/routes"
	"syscall"
)

func init() {
	// Load .env file
	configs.LoadEnv()
	// load keys configuration file
	configs.LoadConfig()
}

func main() {
	log.Printf("Starting server on port %s\n", configs.Env("API_PORT", "8080"))
	r := routes.HandleRequests()

	err := r.Run(":" + configs.Env("API_PORT", "8080"))
	if err != nil {
		log.Fatalf("Couldn't start server due to: %s\n", err.Error())
	}

	// Handle signals like Ctrl+c, Ctrl+ d, etc. with graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 API shutting down...")

}
