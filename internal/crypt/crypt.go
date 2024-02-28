package crypt

import "flag"

type Input struct {
	InPath  string
	OutPath string
}

func NewInput() *Input {
	i := Input{}
	flag.StringVar(&i.InPath, "in", "", "input path")
	flag.StringVar(&i.OutPath, "out", "", "output path")
	flag.Parse()
	return &i
}
