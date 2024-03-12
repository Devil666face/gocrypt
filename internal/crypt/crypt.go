package crypt

import "flag"

type Input struct {
	InPath string
}

func NewInput() *Input {
	i := Input{}
	flag.StringVar(&i.InPath, "d", "", "input direcroty")
	flag.Parse()
	return &i
}
