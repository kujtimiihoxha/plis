package api

import (
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"os"
)

type FsAPI struct {
	fs afero.Fs
}

func (f *FsAPI) ReadFile(path string) (string, error) {
	d, err := afero.ReadFile(f.fs, path)
	return string(d), err
}

func (f *FsAPI) WriteFile(path string, data string) error {

	if b, _ := f.Exists(path); b && !viper.GetBool("plis.fs.force_override") {
		s, _ := f.ReadFile(path)
		if s == data {
			logger.GetLogger().Warnf("`%s` exists and is identical it will be ignored", path)
			return nil
		}
		b := prompter.YN(fmt.Sprintf("`%s` already exists do you want to override it ?", path), false)
		if !b {
			return nil
		}
	}
	return afero.WriteFile(f.fs, path, []byte(data), os.ModePerm)
}

func (f *FsAPI) Mkdir(path string) error {
	return f.fs.Mkdir(path, os.ModePerm)
}

func (f *FsAPI) MkdirAll(path string) error {
	return f.fs.MkdirAll(path, os.ModePerm)
}
func (f *FsAPI) FilePathSeparator() string {
	return afero.FilePathSeparator
}
func (f *FsAPI) Exists(path string) (bool, error) {
	return afero.Exists(f.fs, path)
}
func (f *FsAPI) Walk(root string, fc func(path string, info os.FileInfo, err error) error) error {
	return afero.Walk(f.fs, root, fc)
}
func NewFsAPI(fs afero.Fs) *FsAPI {
	return &FsAPI{
		fs: fs,
	}
}
