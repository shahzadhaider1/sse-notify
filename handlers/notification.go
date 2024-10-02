package handlers

import (
	"bufio"
	"fmt"
	"net/http"
	"sse-notify/notifications"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

type SSEHandler struct {
	nm *notifications.NotificationManager
}

func NewSSEHandler(nm *notifications.NotificationManager) *SSEHandler {
	return &SSEHandler{nm: nm}
}

// ServeSSE handles Server-Sent Events for real-time notifications
func (h *SSEHandler) ServeSSE(c fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return fiber.NewError(http.StatusBadRequest, "UserID is required")
	}
	clientID := time.Now().Format("20060102150405")

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	notificationChannel := h.nm.RegisterClient(userID, clientID)
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		if err := w.Flush(); err != nil {
			fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

			return
		}

		for {
			select {
			case message := <-notificationChannel:
				fmt.Println("got a msg on notificationChannel")
				// fmt.Fprintf(w, "event: notification\n\n data: %s\n\n", message)
				sanitizedMessage := strings.ReplaceAll(message, "\n", "")
				fmt.Fprintf(w, "event: notification\ndata: %s\n\n", sanitizedMessage)

				if err := w.Flush(); err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
					h.nm.UnregisterClient(userID, clientID)
					return
				}
				time.Sleep(5 * time.Second)
			default:
				fmt.Println("sleeping")
				time.Sleep(10 * time.Second)
			}
		}
	}))

	return nil
}
