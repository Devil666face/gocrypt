package io

import (
	"io"
	"log"
	"os"
)

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
