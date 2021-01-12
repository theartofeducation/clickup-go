package clickup

// Status holds the various statuses of a ClickUp Task
type Status string

// ClickUp Statuses
const (
	StatusReadyForDevelopment Status = "ready for development"
	StatusAcceptance          Status = "acceptance"
)
