package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sse-notify/db"
	"sse-notify/handlers"
	"sse-notify/notifications"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

const shutdownTimeout = 5 * time.Second

func main() {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	notificationsManager := notifications.NewNotificationManager(1)
	database := db.NewDB(notificationsManager)
	sseHandler := handlers.NewSSEHandler(notificationsManager)

	app.Get("/notifications/:userID", sseHandler.ServeSSE)
	go func() {
		for i := 0; i < 10; i++ {
			err := database.SaveNotification("user1", "New notification for user1")
			if err != nil {
				fmt.Printf("Error saving notification: %v\n", err)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		if err := app.Listen(":8080"); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	gracefulShutdown(app)
}

func gracefulShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	fmt.Println("\nShutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	}

	select {
	case <-ctx.Done():
		fmt.Println("Server shutdown completed")
	case <-time.After(shutdownTimeout):
		fmt.Println("Server forced to shutdown after timeout")
	}
}
