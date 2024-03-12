package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	s "sync"

	"github.com/Devil666face/gocrypt/internal/async"
	"github.com/Devil666face/gocrypt/internal/crypt"
	"github.com/Devil666face/gocrypt/internal/io"
	"github.com/Devil666face/gocrypt/internal/sync"
	"github.com/Devil666face/gocrypt/pkg/lib"
)

var ErrNotEncrypted = errors.New("not encrypted prefix")

type decrypt struct {
	in   *crypt.Input
	file lib.File
}

func (e *decrypt) to() error {
	if !strings.HasSuffix(e.file.Path, ".enc") {
		return fmt.Errorf("%s: %w", e.file.Path, ErrNotEncrypted)
	}
	in, err := io.ReadFile(e.file.Path)
	if err != nil {
		return err
	}
	chipRsaKey := in[len(in)-988:]
	a := async.New(chipRsaKey, []byte(e.in.InPassPhrase), func(a *async.Async) {
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
	if err := io.WriteFile(plain, strings.TrimSuffix(e.file.Path, ".enc")); err != nil {
		return err
	}
	if err := os.Remove(e.file.Path); err != nil {
		return err
	}
	log.Print(e.file.Path, " decrypted")
	return nil
}

//go:embed id_rsa
var PrivKey string

func worker(jobs <-chan decrypt, results chan<- error, wg *s.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		err := job.to()
		results <- err
	}
}

func main() {
	in := crypt.NewInput()
	// files, err := lib.Walk(in.InPath)
	// if err != nil {
	// 	log.Print(err)
	// }

	files := lib.MustWalk(in.InPath)

	jobs := make(chan decrypt, len(files))
	results := make(chan error, len(files))
	wg := s.WaitGroup{}
	// Start workers
	for i := 1; i <= runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Enqueue jobs
	for _, file := range files {
		e := decrypt{in: in, file: file}
		jobs <- e
	}
	close(jobs)
	wg.Wait()
	close(results)
	for err := range results {
		if err != nil {
			log.Println(err)
		}
	}
}
