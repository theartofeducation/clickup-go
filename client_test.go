package clickup

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

const (
	key    = "abc123"
	secret = "def456"
)

func Test_Client(t *testing.T) {
	t.Run("it creates a new ClickUp Client", func(t *testing.T) {
		client := NewClient(key, secret)

		if _, ok := client.(CUClient); !ok {
			t.Errorf("client is not a CUClient")
		}
	})
}

func Test_VerifySignature(t *testing.T) {
	t.Run("it verifies a valid signature", func(t *testing.T) {
		client := NewClient(key, secret)

		body := []byte("ghi890")

		hash := hmac.New(sha256.New, []byte(secret))
		hash.Write(body)
		signature := hex.EncodeToString(hash.Sum(nil))

		err := client.VerifySignature(signature, body)

		if err != nil {
			t.Errorf("valid signature was not validated")
		}
	})

	t.Run("it does not verify an invalid signature", func(t *testing.T) {
		client := NewClient(key, secret)

		body := []byte("ghi890")

		hash := hmac.New(sha256.New, []byte("bad secret"))
		hash.Write(body)
		signature := hex.EncodeToString(hash.Sum(nil))

		err := client.VerifySignature(signature, body)

		if err != ErrSignatureMismatch {
			t.Errorf("correct error was not returned for invalid signature: got %s want %s", err, ErrSignatureMismatch)
		}
	})
}

func Test_ParseWebhook(t *testing.T) {
	t.Run("it parses and returns a webhook", func(t *testing.T) {
		client := NewClient(key, secret)

		body := ioutil.NopCloser(strings.NewReader(`{"webhook_id": "hgi789", "event": "taskStatusUpdated", "task_id": "test1"}`))

		webhook, err := client.ParseWebhook(body)

		if err != nil {
			t.Fatalf("received error when parsing webhook: %s", err)
		}

		if webhook.ID != "hgi789" {
			t.Errorf("webhook has unexpected ID: got %q want %q", webhook.ID, "hgi789")
		}

		if webhook.Event != EventTaskStatusUpdated {
			t.Errorf("webhook has unexpected event: got %q want %q", webhook.Event, EventTaskStatusUpdated)
		}

		if webhook.TaskID != "test1" {
			t.Errorf("webhook has task ID: got %q want %q", webhook.TaskID, "test1")
		}
	})

	t.Run("it returns an error if webhook cannot be parsed", func(t *testing.T) {
		client := NewClient(key, secret)

		body := ioutil.NopCloser(strings.NewReader(`{"webhook_id": "hgi789", "event": "taskStatusUpdated", "task_id": "test1"`))

		_, err := client.ParseWebhook(body)

		if err == nil {
			t.Fatal("did not receive error when expecting one")
		}

		want := errors.New("Could not parse Webhook body: unexpected EOF")

		if err.Error() != want.Error() {
			t.Errorf("received incorrect error: got %q want %q", err, want)
		}
	})
}

func Test_GetTask(t *testing.T) {
	t.Run("it gets a task", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"id": "test1", "name": "Test 1"}`))
		}))
		defer testServer.Close()

		client := Client{url: testServer.URL}

		task, err := client.GetTask("test1")

		if err != nil {
			t.Fatalf("received error when not expecting one: %s", err)
		}

		if task.ID != "test1" {
			t.Errorf("task has unexpected ID: got %q want %q", task.ID, "test1")
		}

		if task.Name != "Test 1" {
			t.Errorf("task has unexpected name: got %q want %q", task.Name, "Test 1")
		}
	})

	t.Run("it handles HTTP error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusUnauthorized)
		}))
		defer testServer.Close()

		client := Client{url: testServer.URL}

		_, err := client.GetTask("test1")

		if err == nil {
			t.Fatalf("did not receive error when expecting one")
		}

		want := "ClickUp returned status 401"
		if err.Error() != want {
			t.Fatalf("received unexpected error: got %q want %q", err.Error(), want)
		}
	})

	t.Run("it handles JSON parse error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"id": "test1", "name": "Test 1"`))
		}))
		defer testServer.Close()

		client := Client{url: testServer.URL}

		_, err := client.GetTask("test1")

		if err == nil {
			t.Fatalf("did not receive error when expecting one")
		}

		want := "Could not parse Task body: unexpected EOF"
		if err.Error() != want {
			t.Fatalf("received unexpected error: got %q want %q", err.Error(), want)
		}
	})
}
