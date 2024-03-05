package async

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type message []byte

func (m message) messageToPem() []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "MESSAGE", Bytes: m})
}

func (m message) pemToMessage() ([]byte, error) {
	block, _ := pem.Decode(m)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	return block.Bytes, nil
}

type Key []byte

func (k Key) pemToPubRSA() (*rsa.PublicKey, error) {
	block, _ := pem.Decode(k)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func (k Key) pemToPrivRSA() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(k)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

type Async struct {
	text    []byte
	PubKey  Key
	PrivKey Key
}

func New(_text []byte, opts ...func(*Async)) *Async {
	a := &Async{
		text: _text,
	}
	for _, f := range opts {
		f(a)
	}
	return a
}

func (a *Async) Encrypt() ([]byte, error) {
	rsaPubKey, err := a.PubKey.pemToPubRSA()
	if err != nil {
		return nil, fmt.Errorf("error to get pub key: %w", err)
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, a.text, []byte(""))
	if err != nil {
		return nil, fmt.Errorf("error to encrypt message: %w", err)
	}
	return message(ciphertext).messageToPem(), nil
}

func (a *Async) Decrypt() ([]byte, error) {
	rsaPrivKey, err := a.PrivKey.pemToPrivRSA()
	if err != nil {
		return nil, fmt.Errorf("error to get priv key: %w", err)
	}
	pemMess, err := message(a.text).pemToMessage()
	if err != nil {
		return nil, fmt.Errorf("error to convert pem message to bytes: %w", err)
	}
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivKey, pemMess, []byte(""))
	if err != nil {
		return nil, fmt.Errorf("error to decrypt message: %w", err)
	}
	return plaintext, nil
}

func (a *Async) EncryptBase64() ([]byte, error) {
	text, err := a.Encrypt()
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(text)), nil
}

func (a *Async) DecryptBase64() ([]byte, error) {
	text, err := base64.StdEncoding.DecodeString(string(a.text))
	if err != nil {
		return nil, err
	}
	a.text = text
	return a.Decrypt()
}

func GenerateRSAkeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	//nolint:gosec // Error ignore
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func PubRSAtoPEMstr(pubkey *rsa.PublicKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pubkey)})
}

func PrivRSAtoPEMstr(privatekey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privatekey)})
}

func MessageToPEM(msg []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "MESSAGE", Bytes: msg})
}

func PEMtoMessage(pemString string) ([]byte, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	return block.Bytes, nil
}

// func PEMstrToPrivRSA(privPEM []byte) (*rsa.PrivateKey, error) {
// 	block, _ := pem.Decode(privPEM)
// 	if block == nil {
// 		return nil, errors.New("failed to parse PEM block containing the key")
// 	}
// 	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return priv, nil
// }

// func PEMstrToPubRSA(pubPEM []byte) (*rsa.PublicKey, error) {
// 	block, _ := pem.Decode(pubPEM)
// 	if block == nil {
// 		return nil, errors.New("failed to parse PEM block containing the key")
// 	}
// 	return x509.ParsePKCS1PublicKey(block.Bytes)
// }
