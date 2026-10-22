package main

import (
	"api/modules"
	"api/persistence"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load("../.env")

	if err != nil {
		log.Println("Notice: .env file not found, using system environment variables")
	}

	// Connect Database
	db, err := persistence.Connect()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if sqlDB, ok := db.Db.(*sql.DB); ok {
		defer func(sqlDB *sql.DB) {
			err := sqlDB.Close()
			if err != nil {
				log.Printf("Error closing database: %v", err)
			}
		}(sqlDB)
	}

	// Initialize Fiber App
	app := fiber.New(fiber.Config{AppName: "Playground api Server v1.0"})

	//Register Services
	modules.Register(app, db)

	// Save Routes
	if os.Getenv("PRODUCTION") != "true" {
		generateRoutesFile(app)
	}

	// Run The Server
	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	log.SetFlags(0)

	log.Println("\nRunning on port: " + port)
	log.Println("Local: http://localhost:" + port)

	log.Fatal(app.Listen(":" + port))
}

func generateRoutesFile(app *fiber.App) {
	var sb strings.Builder
	sb.WriteString("——— REGISTERED API ROUTES ———\n\n")

	for _, route := range app.GetRoutes() {
		methodFormatted := fmt.Sprintf("[%s]", route.Method)

		line := fmt.Sprintf("%-8s %s\n", methodFormatted, route.Path)
		sb.WriteString(line)
	}

	err := os.WriteFile("routes.txt", []byte(sb.String()), 0644)
	if err != nil {
		log.Printf("Failed to generate routes.txt: %v", err)
		return
	}

	log.Println("Successfully generated routes.txt")
}
