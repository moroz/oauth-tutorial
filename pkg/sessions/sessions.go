package sessions

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"time"
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
	hash.Write(value)
	digest := hash.Sum(nil)
	return digest[:12]
}

func SerializeTimestamp(ts time.Time) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, ts.Unix())
	return buf.Bytes()
}

func ReadTimestamp(src []byte) (result time.Time, err error) {
	var unix int64
	buf := bytes.NewBuffer(src)
	if err = binary.Read(buf, binary.LittleEndian, &unix); err != nil {
		return
	}
	return time.Unix(unix, 0), nil
}

func (s *SessionStorage) Encrypt(value []byte) (nonce, ciphertext []byte, err error) {
	return s.EncryptWithTimestamp(value, time.Now())
}

func (s *SessionStorage) EncryptWithTimestamp(value []byte, timestamp time.Time) (nonce, ciphertext []byte, err error) {
	nonce = DeriveNonce(value, s.deterministicKey[:])
	block, err := aes.NewCipher(s.encryptionKey[:])
	if err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	ts := SerializeTimestamp(timestamp)
	msg := append(ts, value...)

	return nonce, aesgcm.Seal(nil, nonce, msg, nil), nil
}

func (s *SessionStorage) Decrypt(nonce, ciphertext []byte, maxAge int64) (plaintext []byte, err error) {
	block, err := aes.NewCipher(s.encryptionKey[:])
	if err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	decrypted, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("Decrypt: %w", err)
	}
	ts, err := ReadTimestamp(decrypted[:8])
	if err != nil {
		return nil, fmt.Errorf("Decrypt: Invalid timestamp: %w", err)
	}
	age := time.Now().Unix() - ts.Unix()
	if age > maxAge {
		return nil, fmt.Errorf("Message expired")
	}

	return decrypted[8:], nil
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

func (s *SessionStorage) DecodeCookie(cookie string, maxAge int64) (any, error) {
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

	return s.Decrypt(nonce, ciphertext, maxAge)
}
