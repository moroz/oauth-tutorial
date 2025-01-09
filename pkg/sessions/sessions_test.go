package sessions_test

import (
	"crypto/rand"
	"testing"

	"github.com/moroz/oauth-tutorial/pkg/sessions"
)

var storage *sessions.SessionStorage
var encryptionKey [16]byte
var deterministicKey [32]byte

func init() {
	rand.Read(encryptionKey[:])
	rand.Read(deterministicKey[:])
	storage = sessions.NewSessionStorage(encryptionKey, deterministicKey)
}

func TestEncryptDecrypt(t *testing.T) {
	original := "Ich verstehe nur Bahnhof!"
	nonce, ciphertext, err := storage.Encrypt([]byte(original))
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	actual, err := storage.Decrypt(nonce, ciphertext, 3600)
	if string(actual) != original {
		t.Errorf("Expected decrypted value to equal original text %q, got: %q", original, actual)
	}
}
