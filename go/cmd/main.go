package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/google-dev-groups-gmu/ghost/go/internal/api"
	"github.com/google-dev-groups-gmu/ghost/go/internal/firestore"
)

func main() {
	// load env var
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, relying on system env vars")
	}

	// init firestore client
	if err := firestore.Init(); err != nil {
		log.Fatalf("failed to initialize Firestore: %v", err)
	}
	defer firestore.Close()

	// initialize Gin router
	if os.Getenv("DEV") == "false" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()

	FRONTEND_URL := os.Getenv("FRONTEND_URL")
	if FRONTEND_URL == "" {
		log.Printf("frontend url not set in env")
		FRONTEND_URL = "http://localhost:3000"
	}
	config.AllowOrigins = []string{FRONTEND_URL}
	config.AllowCredentials = true
	config.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	r.Use(cors.New(config))

	// health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "health check ok"})
	})

	// root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "GDG ghost map API", "version": "v1.0"})
	})

	// API routes
	a := r.Group("/api")
	{
		// schedule for a specific room
		a.GET("/room", api.GetSpecificRoom)

		// list of rooms and schedules for a specific building
		a.GET("/rooms", api.GetRooms)

		// static building lat/long data
		a.GET("/buildings", api.GetBuildings)
	}

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	r.Run(":" + port)
}
