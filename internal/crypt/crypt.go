package crypt

import "flag"

type Input struct {
	InPath       string
	InPassPhrase string
}

func NewInput() *Input {
	i := Input{}
	flag.StringVar(&i.InPath, "d", "", "input direcroty")
	flag.StringVar(&i.InPassPhrase, "p", "", "passphrase")
	flag.Parse()
	if i.InPath == "" {
		i.InPath = "."
	}
	return &i
}
