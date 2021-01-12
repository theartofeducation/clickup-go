package clickup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

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

// GetTask fetches and returns a Task from ClickUp.
func (c Client) GetTask(taskID string) (Task, error) {
	httpClient := &http.Client{}

	url := c.url + "/task/" + taskID

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("Authorization", c.key)
	request.Header.Add("Content-Type", "application/json")

	var task Task

	response, err := httpClient.Do(request)
	if err != nil {
		return task, errors.Wrap(err, "Could not send request to the ClickUp API")
	}

	if response.StatusCode != http.StatusOK {
		return task, errors.New(fmt.Sprint("ClickUp returned status ", response.StatusCode))
	}

	if err := json.NewDecoder(response.Body).Decode(&task); err != nil {
		return task, errors.Wrap(err, "Could not parse Task body")
	}

	return task, nil
}

// UpdateTaskRequest holds the payload for UpdateTask.
type UpdateTaskRequest struct {
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	Status       Status    `json:"status,omitempty"`
	Priority     int       `json:"priority,omitempty"`
	TimeEstimate int       `json:"time_estimate,omitempty"`
	Assignees    Assignees `json:"assignees,omitempty"`
	Archived     bool      `json:"archived,omitempty"`
}

// Assignees holds lists of users to be added or removed from the Task
type Assignees struct {
	Add []int `json:"add,omitempty"`
	Rem []int `json:"rem,omitempty"`
}

// UpdateTask makes changes to a Task on ClickUp.
// Only pass information that is changed in the task parameter.
func (c Client) UpdateTask(taskID string, task UpdateTaskRequest) error {
	httpClient := &http.Client{}

	url := c.url + "/task/" + taskID

	body, err := json.Marshal(task)
	if err != nil {
		return errors.Wrap(err, "Could not create UpdateTask body")
	}

	request, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	request.Header.Add("Authorization", c.key)
	request.Header.Add("Content-Type", "application/json")

	response, err := httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "Could not send request to the ClickUp API")
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprint("ClickUp returned status ", response.StatusCode))
	}

	var updatedTask Task
	if err := json.NewDecoder(response.Body).Decode(&updatedTask); err != nil {
		return errors.Wrap(err, "Could not parse Task body")
	}

	if updatedTask.Status.Status != task.Status {
		return ErrStatusNotUpdated
	}

	return nil
}
