package sessions

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
)

type SessionStorage struct {
	deterministicKey [32]byte
	encryptionKey    [16]byte
}

func NewSessionStorage(encryptionKey [16]byte, deterministicKey [32]byte) *SessionStorage {
	return &SessionStorage{
		deterministicKey: deterministicKey,
		encryptionKey:    encryptionKey,
	}
}

func DeriveNonce(value, deterministicKey []byte) []byte {
	hash := hmac.New(sha256.New, deterministicKey)
	digest := hash.Sum(value)
	return digest[:12]
}

func (s *SessionStorage) Encrypt(value []byte) (nonce, ciphertext []byte, err error) {
	nonce = DeriveNonce(value, s.deterministicKey[:])
	block, err := aes.NewCipher(s.encryptionKey[:])
	if err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	return nonce, aesgcm.Seal(nil, nonce, value, nil), nil
}

func (s *SessionStorage) Decrypt(nonce, ciphertext []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(s.encryptionKey[:])
	if err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func (s *SessionStorage) EncodeCookie(value any) (string, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(value)
	if err != nil {
		return "", err
	}

	nonce, ciphertext, err := s.Encrypt(buf.Bytes())
	return base64.RawURLEncoding.EncodeToString(nonce) + "." + base64.RawStdEncoding.EncodeToString(ciphertext), nil
}

func (s *SessionStorage) DecodeCookie(cookie string) (any, error) {
	segments := strings.Split(cookie, ".")
	if len(segments) != 2 {
		return nil, errors.New("session cookie must have two segments")
	}

	nonce, err := base64.RawURLEncoding.DecodeString(segments[0])
	if err != nil {
		return nil, fmt.Errorf("DecodeCookie: failed to decode nonce: %w", err)
	}

	ciphertext, err := base64.RawStdEncoding.DecodeString(segments[1])
	if err != nil {
		return nil, fmt.Errorf("DecodeCookie: failed to decode ciphertext: %w", err)
	}

	return s.Decrypt(nonce, ciphertext)
}
