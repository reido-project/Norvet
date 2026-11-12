package main

import (
	"api/libs/nvt"
	"api/libs/nvt/nvtdberr"
	"api/libs/nvt/nvtsb"
	"api/modules"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
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

	// Initialize Norvet context
	diagnosticPtr := flag.Bool("D", false, "Run application in diagnostic mode")
	timePtr := flag.Bool("disable-time", true, "Disable time on terminal log")
	dirPtr := flag.String("dir", "logs/api", "Directory path for log output")
	flag.Parse()

	nvt.Init(nvt.Config{
		Diagnostic:  *diagnosticPtr,
		NvtLogTime:  *timePtr,
		NvtLogColor: true,
		LogOutput:   *dirPtr,
	})

	// Connect Database
	db, err := sqlbuilder.ConnectGoqu()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}

	if sqlDB, ok := db.Db.(*sql.DB); ok {
		defer func(sqlDB *sql.DB) {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing database: %v\n", err)
			}
		}(sqlDB)
	}

	// Initialize Norvet SQL Error Identifier (nvtdb)
	loadCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = nvtdberr.LoadIndexes(loadCtx, db)
	cancel()

	//register Services
	modules.Register(app, db)

	// Save Routes
	if nvt.Ctx() != nil && nvt.Ctx().Diagnostic {
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
