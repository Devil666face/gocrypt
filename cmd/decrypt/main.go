package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	s "sync"

	"github.com/Devil666face/gocrypt/internal/async"
	"github.com/Devil666face/gocrypt/internal/io"
	"github.com/Devil666face/gocrypt/internal/sync"
	"github.com/Devil666face/gocrypt/pkg/lib"
)

var ErrNotEncrypted = errors.New("not encrypted prefix")

func Decrypt(f lib.File) error {
	if !strings.HasSuffix(f.Path, ".enc") {
		return fmt.Errorf("%s: %w", f.Path, ErrNotEncrypted)
	}
	in, err := io.ReadFile(f.Path)
	if err != nil {
		return err
	}
	chipRsaKey := in[len(in)-988:]
	a := async.New(chipRsaKey, func(a *async.Async) {
		a.PrivKey = async.Key(PrivKey)
	})
	aesKey, err := a.DecryptBase64()
	if err != nil {
		return err
	}
	s := sync.New(in[0:len(in)-988], func(s *sync.Sync) {
		s.AesKey = aesKey
	})
	plain, err := s.Decrypt()
	if err != nil {
		return err
	}
	if err := io.WriteFile(plain, strings.TrimSuffix(f.Path, ".enc")); err != nil {
		return err
	}
	if err := os.Remove(f.Path); err != nil {
		return err
	}
	log.Print(f.Path, " decrypted")
	return nil
}

//go:embed id_rsa
var PrivKey string

func main() {
	files, err := lib.Walk(".")
	if err != nil {
		log.Fatal(err)
	}

	results := make(chan error, len(files))
	go func() {
		wg := s.WaitGroup{}
		for _, file := range files {
			wg.Add(1)
			go func(file lib.File) {
				defer wg.Done()
				results <- Decrypt(file)
			}(file)
		}
		wg.Wait()
		close(results)
	}()

	for err := range results {
		if err != nil {
			log.Print(err)
		}
	}
}
