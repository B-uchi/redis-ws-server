package main

import (
	"fmt"
	"huddle-ws-server/database"
	"huddle-ws-server/handler"
	"huddle-ws-server/middleware"
	"huddle-ws-server/rd"
	"huddle-ws-server/ws"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found. Skipping...")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	database.ConnectDatabase()

	go ws.WsManager.Start()

	rd.InitRedis()

	handler.StartRedisListener()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", middleware.WsAuthRequired(), websocket.New(ws.WebsocketHandler))

	fmt.Println("Starting server on port", port)

	app.Listen(":" + port)
}
