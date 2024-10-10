package main

import (
    "context"
    "encoding/json"
    "log"
	"fmt"

    "github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:6379"

// Define the payload struct
type EmailTaskPayload struct {
    UserID    string `json:"userID"`
    Recipient string `json:"recipient"`
    Subject   string `json:"subject"`
    Body      string `json:"body"`
}

// Handler for the "send_email" task
func handleEmailTask(ctx context.Context, t *asynq.Task) error {
    var payload EmailTaskPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        log.Printf("Error unmarshaling payload: %v", err)
        return err
    }

    // Extract payload
    userID := payload.UserID
    recipient := payload.Recipient
    subject := payload.Subject
    body := payload.Body 

    // Simulate sending email
    log.Printf("Sending email with subject %s message %s to user %s at %s", subject, body, userID, recipient)

    // Simulate an error for demonstration purposes
    // Remove this in production
    if userID == "0" {
        err := fmt.Errorf("invalid user ID: %s", userID)
        log.Printf("Error sending email: %v", err)
        return err
    }
	
	// This line is **removed** to avoid logging twice:
	// log.Printf("Sending email to %s (UserID: %s)\n", recipient, userID)
	
	return nil
}

func main() {
    // Create a new Asynq server with Redis as the broker
    server := asynq.NewServer(
        asynq.RedisClientOpt{Addr: redisAddr},
        asynq.Config{
            Concurrency: 10, // Number of workers
        },
    )

    // Register the task handler
    mux := asynq.NewServeMux()
    mux.HandleFunc("send_email", handleEmailTask)

    // Start processing tasks
    if err := server.Run(mux); err != nil {
        log.Fatalf("Could not run server: %v", err)
    }
}