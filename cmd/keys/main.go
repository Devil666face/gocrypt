package main

import (
	"log"
	"os"

	"github.com/Devil666face/gocrypt/internal/async"
	"github.com/Devil666face/gocrypt/internal/io"
)

func main() {
	priv, pub := async.GenerateRSAkeyPair()
	if err := io.WriteFile(async.PrivRSAtoPEMstr(priv), "id_rsa"); err != nil {
		log.Fatal(err)
	}
	if err := io.WriteFile(async.PubRSAtoPEMstr(pub), "id_rsa.pub"); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)

}
