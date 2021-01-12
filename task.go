package clickup

import (
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
