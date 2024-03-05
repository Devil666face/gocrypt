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

//go:embed id_rsa
var PrivKey string

func main() {
	i := crypt.NewInput()
	in, err := io.ReadFile(i.InPath)
	if err != nil {
		log.Fatal(err)
	}
	chipRsaKey := in[len(in)-988:]
	a := async.New(chipRsaKey, func(a *async.Async) {
		a.PrivKey = async.Key(PrivKey)
	})
	aesKey, err := a.DecryptBase64()
	if err != nil {
		log.Fatal(err)
	}
	s := sync.New(in[0:len(in)-988], func(s *sync.Sync) {
		s.AesKey = aesKey
	})
	plain, err := s.Decrypt()
	if err != nil {
		log.Fatal(err)
	}
	if err := io.WriteFile(plain, i.OutPath); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
