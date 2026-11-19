package main

import (
	"api/modules"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"gitlab.com/Norvet/nvtdberr"
	sqlbuilder "gitlab.com/Norvet/nvtsb"
	"gitlab.com/Norvet/nvtsql"
)

func main() {
	// Load .env
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Fatal: .env file not found. Please copy .env.example to .env\nGuide/Template: https://github.com/reido-project/Norvet\nDetails: %v", err)
	}

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("%s v%s", os.Getenv("APP_NAME"), os.Getenv("APP_VERSION")),
	})

	// Initialize Norvet libraries
	diagnostic, err := initNorvet()
	if err != nil {
		log.Fatalf("Failed to initialize Norvet: %v\n", err)
	}

	// Connect Database
	db, err := sqlbuilder.ConnectGoqu()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}

	defer func() {
		if err := nvtsql.CloseAll(); err != nil {
			log.Printf("Error closing database: %v\n", err)
		}
	}()

	// Initialize Norvet SQL Error Identifier (nvtdb)
	loadCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = nvtdberr.LoadIndexes(loadCtx, db)
	cancel()

	//register Services
	modules.Register(app, db)

	// Save Routes
	if diagnostic {
		generateRoutesFile(app)
	}

	// Run The Server
	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	log.SetFlags(0)

	log.Println("Running on port: " + port)
	log.Println("Local: http://localhost:" + port + "\n")

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
		log.Printf("Failed to generate routes.txt: %v\n", err)
		return
	}

	log.Println("Successfully generated routes.txt")
}
