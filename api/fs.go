package api

import (
	"github.com/spf13/afero"
	"os"
	"strings"
)

type FsAPI struct {
	fs afero.Fs
}

func (f *FsAPI) ReadFile(path string) (string, error) {
	path = strings.Replace(path,"\\",f.FileSeparator(),-1)
	d, err := afero.ReadFile(f.fs, path)
	return string(d), err
}

func (f *FsAPI) WriteFile(path string, data string) error {
	path = strings.Replace(path,"\\",f.FileSeparator(),-1)
	return afero.WriteFile(f.fs, path, []byte(data), os.ModePerm)
}

func (f *FsAPI) Mkdir(path string) error {
	path = strings.Replace(path,"\\",f.FileSeparator(),-1)
	return f.fs.Mkdir(path, os.ModePerm)
}

func (f *FsAPI) MkdirAll(path string) error {
	path = strings.Replace(path,"\\",f.FileSeparator(),-1)
	return f.fs.MkdirAll(path, os.ModePerm)
}
func (f *FsAPI) FileSeparator() string {
	return afero.FilePathSeparator
}
func (f *FsAPI) Exists(path string) (bool, error) {
	path = strings.Replace(path,"\\",f.FileSeparator(),-1)
	return afero.Exists( f.fs,path)
}
func NewFsAPI(fs afero.Fs) *FsAPI {
	return &FsAPI{
		fs: fs,
	}
}
