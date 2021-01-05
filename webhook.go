package clickup

// Event is the action that triggered the Webhook.
type Event string

// Webhook Events.
const (
	EventTaskStatusUpdated = "taskStatusUpdated"
)

// Webhook holds the information for a ClickUp Webhook.
type Webhook struct {
	ID     string `json:"webhook_id"`
	Event  Event
	TaskID string `json:"task_id"`
}
