package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"io"
	"log"
	"os"
)

var (
	// We're using a 32 byte long secret key.
	// This is probably something you generate first
	// then put into and environment variable.
	secretKey string = "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
)

func encrypt(plain []byte) ([]byte, error) {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	return gcm.Seal(nonce, nonce, plain, nil), nil
}

func decrypt(ciphertext []byte) ([]byte, error) {
	aes, err := aes.NewCipher([]byte(secretKey))
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
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	return io.ReadAll(file)
}

func WriteFile(b []byte, path string) error {
	return os.WriteFile(path, b, 0644)
}

func main() {
	var (
		inPath  string
		outPath string
	)
	isEncyprt := flag.Bool("encrypt", false, "encrypt plain data")
	isDecrypt := flag.Bool("decrypt", false, "decrypt data")
	flag.StringVar(&inPath, "in", "", "input path")
	flag.StringVar(&outPath, "out", "", "output path")
	flag.Parse()

	in, err := ReadFile(inPath)
	if err != nil {
		log.Fatal(err)
	}
	if *isEncyprt {
		chip, err := encrypt(in)
		if err != nil {
			log.Fatal(err)
		}
		if err := WriteFile(chip, outPath); err != nil {
			log.Fatal(err)
		}
	}
	if *isDecrypt {
		plain, err := decrypt(in)
		if err != nil {
			log.Fatal(err)
		}
		if err := WriteFile(plain, outPath); err != nil {
			log.Fatal(err)
		}
	}
}
