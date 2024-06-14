package main

import (
	"fiber-search-engine/db"
	"fiber-search-engine/routes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
)

func main() {
	env := godotenv.Load()
	if env != nil {
		panic("Cannot find environment variable")

	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":4000"
	} else {
		port = ":" + port
	}

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	routes.SetRoutes(app)

	app.Use(compress.New())

	db.InitDB()

	// Start out server and listen for a shutdown

	go func() {
		if err := app.Listen(port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // Block main thread until interrupted

	app.Shutdown()
	fmt.Println(("Shutting down server"))
}
