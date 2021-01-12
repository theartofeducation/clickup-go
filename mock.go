package clickup

import (
	"io"

	"github.com/pkg/errors"
)

// ErrTest is returned for testing method errors.
var ErrTest = errors.New("Test error")

// MockClient is a mock Client to use for testing.
type MockClient struct {
	Task                 Task
	VerifySignatureError bool
	ParseWebhookError    bool
	GetTaskError         bool
	UpdateTaskError      bool
}

// GetTask mock fetches and returns a Task from ClickUp.
func (c MockClient) GetTask(taskID string) (Task, error) {
	if c.GetTaskError {
		return Task{}, ErrTest
	}

	return c.Task, nil
}

// ParseWebhook mock parses a Webhook's body and returns a Webhook struct.
func (c MockClient) ParseWebhook(body io.ReadCloser) (Webhook, error) {
	if c.ParseWebhookError {
		return Webhook{}, errors.Wrap(ErrTest, "Could not parse Webhook body")
	}

	return Webhook{}, nil
}

// VerifySignature mock validates a Webhook's signature.
func (c MockClient) VerifySignature(signature string, body []byte) error {
	if c.VerifySignatureError {
		return ErrSignatureMismatch
	}

	return nil
}

// UpdateTask mocks makes changes to a Task on ClickUp.
func (c MockClient) UpdateTask(taskID string, task UpdateTaskRequest) error {
	if c.UpdateTaskError {
		return ErrTest
	}

	return nil
}
