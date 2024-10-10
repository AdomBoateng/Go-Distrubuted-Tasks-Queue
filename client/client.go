package main

import(
	"encoding/json"
	"log"
	"time"
	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:6379"

// **Use the same payload struct name as in the worker**
type EmailTaskPayload struct {  // **changed from MailPayload**
	UserID    string `json:"userID"`
    Recipient string `json:"recipient"`
    Subject   string `json:"subject"`
    Body      string `json:"body"`
}

func main() {
	// Create a new client
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// Create a new task payload (for sending mail)
	payload := EmailTaskPayload{  // **use the same struct as in worker**
		UserID: "001",
		Recipient: "gboateng@quantumgroupgh.com",
		Subject:   "Hello Quantum Group",
        Body:      "This is a test mail from asynq",
    }

	// Serialize payload into a JSON
	taskPayload, err := json.Marshal(payload)
	if err != nil{
		log.Fatalf("Could not marshal payload : %v", err)
	}

	// Create a task and enqueue
	task := asynq.NewTask("send_email", taskPayload)  // **changed "send_mail" to "send_email"**
	info, err := client.Enqueue(task, asynq.MaxRetry(6), asynq.ProcessIn(10*time.Second))
	if err != nil {
		log.Fatalf("Could not enqueue task : %v", err)
	}
	log.Printf("Task enqueued: id=%s queue=%s", info.ID, info.Queue)
}
