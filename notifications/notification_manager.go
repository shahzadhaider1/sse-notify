package notifications

import (
	"fmt"
	"sync"
)

type NotificationManager struct {
	userChannels map[string]map[string]chan string
	mu           sync.RWMutex
	bufferSize   int
}

func NewNotificationManager(bufferSize int) *NotificationManager {
	return &NotificationManager{
		userChannels: make(map[string]map[string]chan string),
		bufferSize:   bufferSize,
	}
}

func (m *NotificationManager) RegisterClient(userID, clientID string) chan string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.userChannels[userID]; !ok {
		m.userChannels[userID] = make(map[string]chan string)
	}

	ch := make(chan string, m.bufferSize)
	m.userChannels[userID][clientID] = ch

	fmt.Printf("Registered client %s for user %s\n", clientID, userID)
	return ch
}

func (m *NotificationManager) UnregisterClient(userID, clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if clients, ok := m.userChannels[userID]; ok {
		if ch, exists := clients[clientID]; exists {
			close(ch)
			delete(clients, clientID)
		}

		if len(clients) == 0 {
			delete(m.userChannels, userID)
		}
	}

	fmt.Printf("Unregistered client %s for user %s\n", clientID, userID)
}

func (m *NotificationManager) SendNotification(userID, notification string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fmt.Println("Sending notification to user : ", userID)
	clients, ok := m.userChannels[userID]
	if !ok {
		fmt.Println("Returning. No clients found for user : ", userID)
		return
	}

	fmt.Println("here are the client channels found : ", clients)

	for clientID, ch := range clients {
		select {
		case ch <- notification:
			fmt.Printf("Notification sent to user %s, client %s\n", userID, clientID)
		default:
			fmt.Printf("Notification dropped for user %s, client %s (channel full)\n", userID, clientID)
		}
	}
}
