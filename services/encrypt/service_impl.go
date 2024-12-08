package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/melvinodsa/go-iam/sdk"
)

type service struct {
	cipher cipher.AEAD
}

func NewService(key sdk.MaskedBytes) (Service, error) {
	keyB := []byte(key)
	bl, err := aes.NewCipher(keyB)
	if err != nil {
		return nil, fmt.Errorf("error creating block: %w", err)
	}
	gcm, err := cipher.NewGCM(bl)
	if err != nil {
		return nil, fmt.Errorf("error creating gcm: %w", err)
	}
	return &service{
		cipher: gcm,
	}, nil
}

func (s service) Encrypt(rawMessage string) (string, error) {
	nonce := make([]byte, s.cipher.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error creating nonce: %w", err)
	}
	ciphertext := s.cipher.Seal(nonce, nonce, []byte(rawMessage), nil)
	enc := hex.EncodeToString(ciphertext)
	return enc, nil
}

func (s service) Decrypt(encryptedMessage string) (string, error) {
	by, err := hex.DecodeString(encryptedMessage)
	if err != nil {
		return "", fmt.Errorf("error decoding encrypted message: %w", err)
	}
	nonceSuze := s.cipher.NonceSize()
	decryptedData, err := s.cipher.Open(nil, by[:nonceSuze], by[nonceSuze:], nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting message: %w", err)
	}
	return string(decryptedData), nil
}
