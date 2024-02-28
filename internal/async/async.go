package async

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func Encrypt(message, pubKey []byte) ([]byte, error) {
	rsaPubKey, err := PEMstrToPubRSA(pubKey)
	if err != nil {
		return nil, fmt.Errorf("error to get pub key: %w", err)
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, message, []byte(""))
	if err != nil {
		return nil, fmt.Errorf("error to encrypt message: %w", err)
	}
	return MessageToPEM(ciphertext), nil
}

func Decrypt(ciphertext, privKey string) (string, error) {
	rsaPrivKey, err := PEMstrToPrivRSA(privKey)
	if err != nil {
		return "", fmt.Errorf("error to get priv key: %w", err)
	}
	// fmt.Println(lib.PrivRSAtoPEMstr(rsaPrivKey))
	pemMess, err := PEMtoMessage(ciphertext)
	if err != nil {
		return "", fmt.Errorf("error to convert pem message to bytes: %w", err)
	}
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivKey, pemMess, []byte(""))
	if err != nil {
		return "", fmt.Errorf("error to decrypt message: %w", err)
	}
	return string(plaintext), nil
}

func GenerateRSAkeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	//nolint:gosec // Error ignore
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func PEMstrToPrivRSA(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func PEMstrToPubRSA(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func PubRSAtoPEMstr(pubkey *rsa.PublicKey) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
}

func PrivRSAtoPEMstr(privatekey *rsa.PrivateKey) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}))
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
