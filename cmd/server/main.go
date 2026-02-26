package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yuyudeqiu/chronicle/internal/handler"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

func main() {
	// Initialize database
	service.InitDB("data/app.db")

	// Setup router
	r := gin.Default()

	// Register routes
	handler.RegisterRoutes(r)

	// Serve static files
	r.StaticFile("/", "./static/index.html")
	r.StaticFile("/app.js", "./static/app.js")
	r.StaticFile("/summary.html", "static/summary.html")

	// Start server
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
