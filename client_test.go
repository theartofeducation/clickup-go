package clickup

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
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
