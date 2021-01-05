package clickup

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// CUClient defines available methods.
type CUClient interface {
	VerifySignature(signature string, body []byte) error
	ParseWebhook(body io.ReadCloser) (Webhook, error)
	GetTask(taskID string) (Task, error)
}

const apiURL = "https://api.clickup.com/api/v2"

// Options are the settings needed when creating a new Client.
type Options struct {
	Key                     string
	TaskStatusUpdatedSecret string
}

// Client handles interaction with the ClickUp API.
type Client struct {
	key                     string
	taskStatusUpdatedSecret string
	url                     string
}

// NewClient creates and returns a new ClickUp Client.
func NewClient(options Options) CUClient {
	client := Client{
		key:                     options.Key,
		taskStatusUpdatedSecret: options.TaskStatusUpdatedSecret,
		url:                     apiURL,
	}

	return client
}

// VerifySignature validates a Webhook's signature.
func (c Client) VerifySignature(signature string, body []byte) error {
	secret := []byte(c.taskStatusUpdatedSecret)

	hash := hmac.New(sha256.New, secret)
	hash.Write(body)
	generatedSignature := hex.EncodeToString(hash.Sum(nil))

	if signature == generatedSignature {
		return nil
	}

	return ErrSignatureMismatch
}

// ParseWebhook parses a Webhook's body and returns a Webhook struct.
func (c Client) ParseWebhook(body io.ReadCloser) (Webhook, error) {
	defer body.Close()

	var webhook Webhook

	if err := json.NewDecoder(body).Decode(&webhook); err != nil {
		return webhook, errors.Wrap(err, "Could not parse Webhook body")
	}

	return webhook, nil
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
