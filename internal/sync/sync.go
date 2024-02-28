package sync

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/Devil666face/gocrypt/pkg/lib"
)

type Sync struct {
	plain  []byte
	AesKey []byte
}

func New(_plain []byte) *Sync {
	return &Sync{
		plain: _plain,
	}
}

func (s *Sync) Encrypt() ([]byte, error) {
	s.AesKey = lib.AES32RandomKey()
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
	return gcm.Seal(nonce, nonce, s.plain, nil), nil
}

// func DecryptAES(ciphertext []byte, secretKey string) ([]byte, error) {
// 	aes, err := aes.NewCipher([]byte(secretKey))
// 	if err != nil {
// 		return nil, err
// 	}

// 	gcm, err := cipher.NewGCM(aes)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Since we know the ciphertext is actually nonce+ciphertext
// 	// And len(nonce) == NonceSize(). We can separate the two.
// 	nonceSize := gcm.NonceSize()
// 	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

// 	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return plain, nil
// }
