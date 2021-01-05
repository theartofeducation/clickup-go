package clickup

// Task holds information for a ClickUp Task.
type Task struct {
	ID     string
	Name   string
	Status TaskStatus
	URL    string
}

// TaskStatus holds a Task's Status information.
type TaskStatus struct {
	Status Status
}
