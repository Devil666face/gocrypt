package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/Devil666face/gocrypt/internal/async"
	"github.com/Devil666face/gocrypt/internal/crypt"
	"github.com/Devil666face/gocrypt/internal/io"
	"github.com/Devil666face/gocrypt/internal/sync"
)

//go:embed id_rsa.pub
var PubKey string

func main() {
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
	a := async.New(s.AesKey, func(a *async.Async) {
		a.PubKey = async.Key(PubKey)
	})
	chipAesKey, err := a.EncryptBase64()
	if err != nil {
		log.Fatal(err)
	}
	chip = append(chip, chipAesKey...)
	if err := io.WriteFile(chip, i.OutPath); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
