package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/Devil666face/gocrypt/internal/async"
	"github.com/Devil666face/gocrypt/internal/crypt"
	"github.com/Devil666face/gocrypt/internal/io"
	"github.com/Devil666face/gocrypt/internal/sync"
)

func keys() {
	priv, pub := async.GenerateRSAkeyPair()
	fmt.Println(async.PrivRSAtoPEMstr(priv))
	fmt.Println(async.PubRSAtoPEMstr(pub))
}

const PubKey = `-----BEGIN RSA PUBLIC KEY-----
MIICCgKCAgEAw3VX83Uoxjtpx2zE4McFtOGq6IYSwc0Dx8fym1yS8ChXLlL6L032
MZlFg9RRbRT0cJLOPFsy9Q9mUKDNoLP3gHvaVIAhlJwMfK10EH3qvnXreqScfOiq
lSHh6rU72xmMiBwNWhlxrmpcxGRqe9q7vOX3KNd2gQjjQwVmFoJyVDp65l/a5JzV
4rxbI7ZcCpPX5O9gPtc98zfNaUpD5R9JNMLqdk/hI/M00VKckJWWiI/tj1ktmFuY
PNuYYSanLK8CQ21RUdDFAX7EM1+hRUDfXoHpHu66Cj8mbDu2XE6Jujo5DhqAiN9t
pj5HA1OVujpVAn4Oq8Sx5yvCc/bzmZM0bMyk1erIFtHWa97klbWSP0TRGv+jaWtJ
nz+ZQzuJR80dzW1bnFQirqaOV2o1nJ9ybHTP4GPe6moUCSOs0KZFK9WEp3ttJYJW
X24eh4ajSGp4Otbz+hCQzewi1W1WqIeyGXBL0pe0G/BGM3yiPX3PDiqyrkKRiFH4
RUCbbA5cL27arvckeQ+blQdHCKA8Js1PHKF9XS+5SQsEcR6qfLojsgfT2JLN8TO0
GHKhDn+E9GfCeVj594Yn6ScBdai/S4IuX3HuM0rawpXVZBu0IvVjqqK4/MWlbynl
NgAUiTWk/7McjOeb04eEL0chsj5sMDbp+xLDXoLbB+/8Tsmi9xPNpr0CAwEAAQ==
-----END RSA PUBLIC KEY-----`

func main() {
	// keys()
	i := crypt.NewInput()
	in, err := io.ReadFile(i.InPath)
	if err != nil {
		log.Fatal(err)

	}
	s := sync.New(in)
	chip, err := s.Encrypt()
	if err != nil {
		log.Fatal(err)
	}
	chipAesKey, err := async.Encrypt(s.AesKey, []byte(PubKey))
	if err != nil {
		log.Fatal(err)
	}
	key := base64.StdEncoding.EncodeToString(chipAesKey)
	chip = append(chip, []byte(key)...)
	if err := io.WriteFile(chip, i.OutPath); err != nil {
		log.Fatal(err)
	}
}
