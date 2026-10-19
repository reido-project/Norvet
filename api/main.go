package main

import (
	"api/modules/user"
	"api/persistence"
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
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
		defer sqlDB.Close()
	}

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Playground api Server v1.0",
	})

	//Register Services
	user.Register(app, db)

	// Run The Server
	port := os.Getenv("APP_PORT")
	log.Println("Server running on port :" + port)
	log.Fatal(app.Listen(":" + port))
}
