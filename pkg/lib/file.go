package lib

import (
	"io/fs"
	"path/filepath"
	"sort"
)

type File struct {
	Path string
	Info fs.FileInfo
}

type walker struct {
	files []File
}

func Walk(path string) ([]File, error) {
	f := walker{}
	if err := filepath.WalkDir(path, f.walk); err != nil {
		return nil, err
	}
	f.sizeSort()
	return f.files, nil
}

func (f *walker) sizeSort() {
	sort.Slice(f.files, func(i, j int) bool {
		return f.files[i].Info.Size() < f.files[j].Info.Size()
	})
}

func (f *walker) walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
		return nil
	}
	i, err := d.Info()
	if err != nil {
		return err
	}
	f.files = append(f.files, File{
		Path: s,
		Info: i,
	})
	return nil
}
