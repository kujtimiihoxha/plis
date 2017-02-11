package api

import (
	"github.com/spf13/afero"
	"os"
)

type FsApi struct {
	fs afero.Fs
}

func (f *FsApi) ReadFile(path string) (string, error) {
	d, err := afero.ReadFile(f.fs, path)
	return string(d), err
}

func (f *FsApi) WriteFile(path string, data string) error {
	return afero.WriteFile(f.fs, path, []byte(data), os.ModePerm)
}

func (f *FsApi) Mkdir(path string) error {
	return f.fs.Mkdir(path, os.ModePerm)
}

func (f *FsApi) MkdirAll(path string) error {
	return f.fs.MkdirAll(path, os.ModePerm)
}
func NewFsApi(fs afero.Fs) *FsApi {
	return &FsApi{
		fs: fs,
	}
}
