package main

import (
	"sse-notify/db"
	"sse-notify/handlers"
	"sse-notify/notifications"
	"time"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	notificationsManager := notifications.NewNotificationManager(1)

	database := db.NewDB(notificationsManager)

	sseHandler := handlers.NewSSEHandler(notificationsManager)

	app.Get("/notifications/:userID", sseHandler.ServeSSE)

	go func() {
		for i := 0; i < 10; i++ {
			_ = database.SaveNotification("user1", "New notification for user1")
			time.Sleep(time.Second * 10)
		}
	}()

	// Start the server
	app.Listen(":8080")
}
