package pii

// https://dev.to/breda/secret-key-encryption-with-go-using-aes-316d

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

const keySize = 32

type PII struct {
	gcm cipher.AEAD
}

func Create(key []byte) PII {
	aes, err := aes.NewCipher(key)
	// Errors can only be KeySizeError which is created because of an incorrect key size
	// 16, 24, 32 are the only sizes accepted. This should never happen.
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	return PII{gcm: gcm}
}

func NewKey() []byte {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	return key
}

func (pii PII) Encrypt(text string) string {
	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, pii.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	// ciperBytes here is actually nonce+ciperBytes
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciperBytes.
	ciperBytes := pii.gcm.Seal(nonce, nonce, []byte(text), nil)

	return base64.StdEncoding.EncodeToString(ciperBytes)
}

func (pii PII) Decrypt(ciphertext string) string {
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		panic(err)
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := pii.gcm.NonceSize()
	nonce, cipherValue := cipherBytes[:nonceSize], cipherBytes[nonceSize:]

	text, err := pii.gcm.Open(nil, nonce, cipherValue, nil)
	if err != nil {
		panic(err)
	}

	return string(text)
}
