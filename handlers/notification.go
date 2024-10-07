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
		fmt.Fprintf(w, ": connected\n\n")
		if err := w.Flush(); err != nil {
			fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
			return
		}

		for {
			select {
			case message := <-notificationChannel:
				sanitizedMessage := strings.ReplaceAll(message, "\n", "")
				fmt.Fprintf(w, "event: notification\ndata: %s\n\n", sanitizedMessage)

				if err := w.Flush(); err != nil {
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
					h.nm.UnregisterClient(userID, clientID)
					return
				}
			default:
				fmt.Fprintf(w, ": keep-alive\n\n")
				if err := w.Flush(); err != nil {
					fmt.Printf("Error while flushing (keep-alive): %v. Closing http connection.\n", err)
					h.nm.UnregisterClient(userID, clientID)
					return
				}
				time.Sleep(10 * time.Second)
			}
		}
	}))

	return nil
}
