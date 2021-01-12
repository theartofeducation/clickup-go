package clickup

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func Test_UpdateTask(t *testing.T) {
	t.Run("it updates a task", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"id": "test1", "name": "Test 1", "status": {"status": "acceptance"}}`))
		}))
		defer testServer.Close()

		client := Client{url: testServer.URL}

		task := UpdateTaskRequest{Status: StatusAcceptance}

		err := client.UpdateTask("test1", task)

		if err != nil {
			t.Fatalf("received error when not expecting one: %s", err)
		}
	})

	t.Run("it returns error for non-ok response", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusUnauthorized)
		}))
		defer testServer.Close()

		client := Client{url: testServer.URL}

		task := UpdateTaskRequest{Status: StatusAcceptance}

		err := client.UpdateTask("test1", task)

		if err == nil {
			t.Fatal("did not receive error when expecting one")
		}

		want := "ClickUp returned status 401"
		if err.Error() != want {
			t.Fatalf("received unexpected error: got %q want %q", err.Error(), want)
		}
	})

	t.Run("it handles json parse error response", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"id": "test1", "name": "Test 1", "status": {"status": "acceptance"`))
		}))
		defer testServer.Close()

		client := Client{url: testServer.URL}

		task := UpdateTaskRequest{Status: StatusAcceptance}

		err := client.UpdateTask("test1", task)

		if err == nil {
			t.Fatal("did not receive error when expecting one")
		}

		want := "Could not parse Task body: unexpected EOF"
		if err.Error() != want {
			t.Fatalf("received unexpected error: got %q want %q", err.Error(), want)
		}
	})

	t.Run("it returns error if status was not updated", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"id": "test1", "name": "Test 1", "status": {"status": "ready for development"}}`))
		}))
		defer testServer.Close()

		client := Client{url: testServer.URL}

		task := UpdateTaskRequest{Status: StatusAcceptance}

		err := client.UpdateTask("test1", task)

		if err == nil {
			t.Fatal("did not receive error when expecting one")
		}

		want := "Task status was not updated"
		if err.Error() != want {
			t.Fatalf("received unexpected error: got %q want %q", err.Error(), want)
		}
	})
}
