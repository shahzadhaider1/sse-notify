# sse-notify
Real time notifications through SSE. This repository implements a Server-Sent Events (SSE) notification system using Golang and Fiber. The system allows clients to subscribe to real-time notifications. Multiple clients can connect to the server and listen for events, such as notifications, which are pushed in real-time without requiring the client to poll the server continuously.

## Overview
The SSE notification system consists of two main components:

### NotificationManager
This manages clients and notifications, ensuring that each connected client for a specific user receives notifications through a channel.

### SSEHandler
This handles the HTTP request from clients who subscribe to the SSE stream and sends notifications over the stream in real time.

## How it works

### Registering clients
When a client connects to the `/notifications/:userID` endpoint, the SSEHandler registers that client under a specific userID and generates a unique clientID using the current timestamp.

Each user can have multiple connected clients, represented by a map of channels where the outer key is the userID and the inner key is the clientID.

### Listening for notifications
Once the client is registered, it starts listening for notifications sent over the channel. The server keeps the connection alive using the SSE protocol, which automatically handles reconnections in case of network issues.

### Sending notifications 
The server (or any part of the system) can send a notification to a specific user by invoking the `SendNotification` method of `NotificationManager`. If the user has active clients, each client will receive the notification through its dedicated channel.

### Unregistering clients
If a client disconnects or an error occurs (e.g., a connection issue), the client will be unregistered from the NotificationManager, and its associated channel will be closed to release resources.


## Key Features

### Real-time notifications
Clients receive notifications as soon as they are sent, without needing to poll the server.

### Scalable
Multiple clients can connect for a single user, and the system ensures each client receives messages.

### Channel-based communication
Each client has its own buffered channel, allowing the system to queue up messages and handle slow consumers.

### Graceful error handling
The system unregisters clients if an error occurs during communication, freeing up resources.


## Future Enhancements

While this SSE-based notification system is functional, there are several ways to improve its scalability, performance, and robustness for production use:

### Persistent Connections & Message Durability
Issue: Currently, notifications are only delivered to active clients, and if a client disconnects, it will miss any notifications sent during the downtime.
Improvement: Introduce message persistence using a message queue (e.g., Kafka or RabbitMQ) to ensure that even disconnected clients can retrieve notifications when they reconnect.

### Horizontal Scalability
Issue: The current implementation is limited to a single instance of the server. If the server crashes or restarts, all clients are disconnected, and no new notifications can be processed.
Improvement: Implement distributed pub/sub systems (e.g., Redis Pub/Sub) or a dedicated message broker (like NATS or Kafka) for notifications, allowing multiple server instances to handle different clients. This would enable load balancing and ensure scalability across distributed systems.

### Efficient Channel Management:
Issue: The system keeps channels open for clients until they are explicitly unregistered or disconnected. This can lead to memory leaks if channels are not closed properly.
Improvement: Implement heartbeat mechanisms or connection timeouts to detect dead clients and close their channels automatically.

### Security & Authentication:
Issue: Currently, there is no authentication mechanism for registering clients. Any client can subscribe to notifications for any userID.
Improvement: Add authentication and authorization layers (e.g., JWT-based authentication) to ensure that only authorized clients can subscribe to a specific user’s notifications.

### Backpressure Management:
Issue: If clients are slow in processing notifications, they might drop messages as channels are buffered but not infinitely large.
Improvement: Implement backpressure mechanisms to handle slow consumers more gracefully, such as by pausing notification delivery or scaling buffers dynamically.

### Performance Metrics:
Improvement: Introduce logging and monitoring tools (such as Prometheus, Grafana) to track the number of active connections, notifications sent, and system performance under load.

# Contributions
Contributions to this project are welcome! If you would like to contribute, please feel free to:

1. Submit pull requests for bug fixes, enhancements, or optimizations.
2. Propose new features by opening issues or discussing your ideas.
3. Improve documentation or add examples to help others get started.
4. Before contributing, please review the contribution guidelines. Let’s work together to make this notification system even better!
