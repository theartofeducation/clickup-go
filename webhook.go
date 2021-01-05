package clickup

// Webhook holds the information for a ClickUp Webhook.
type Webhook struct {
	ID     string `json:"webhook_id"`
	Event  string
	TaskID string `json:"task_id"`
}
