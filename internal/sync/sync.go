package sync

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/Devil666face/gocrypt/pkg/lib"
)

type Sync struct {
	text   []byte
	AesKey []byte
}

func New(_text []byte, opts ...func(*Sync)) *Sync {
	s := &Sync{
		text: _text,
	}
	s.AesKey = lib.AES32RandomKey()
	for _, f := range opts {
		f(s)
	}
	return s
}

func (s *Sync) Encrypt() ([]byte, error) {
	aes, err := aes.NewCipher(s.AesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, s.text, nil), nil
}

func (s *Sync) Decrypt() ([]byte, error) {
	aes, err := aes.NewCipher(s.AesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := s.text[:nonceSize], s.text[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
