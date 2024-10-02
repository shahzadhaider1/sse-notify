package db

import (
	"sse-notify/notifications"
)

type DB struct {
	nm *notifications.NotificationManager
}

func NewDB(nm *notifications.NotificationManager) *DB {
	return &DB{nm: nm}
}

func (db *DB) SaveNotification(userID, notification string) error {
	// Your database logic for saving notification goes here
	// Assume we have saved the notification successfully...

	// Notify subscribed clients
	db.nm.SendNotification(userID, notification)

	return nil
}
