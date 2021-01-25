package clickup

// Status holds the various statuses of a ClickUp Task
type Status string

// ClickUp Statuses
const (
	StatusAcceptance             Status = "acceptance"
	StatusInDevelopmentClubhouse Status = "in development (clubhouse)"
	StatusReadyForDevelopment    Status = "ready for development"
)
