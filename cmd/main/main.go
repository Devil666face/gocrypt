package main

import (
	"flag"
	"log"

	"github.com/Devil666face/gocrypt/internal/io"
	"github.com/Devil666face/gocrypt/internal/sync"
)

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

	in, err := io.ReadFile(inPath)
	if err != nil {
		log.Fatal(err)
	}
	if *isEncyprt {
		chip, err := sync.EncryptAES(in)
		if err != nil {
			log.Fatal(err)
		}
		if err := io.WriteFile(chip, outPath); err != nil {
			log.Fatal(err)
		}
	}
	if *isDecrypt {
		// plain, err := sync.DecryptAES(in)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// if err := io.WriteFile(plain, outPath); err != nil {
		// 	log.Fatal(err)
		// }
	}
}
