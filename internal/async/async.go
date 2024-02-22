package async

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func PubRSAtoPEMstr(pubkey *rsa.PublicKey) string {
	return string(pem.EncodeToMemory(&pem.Block{Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
}
