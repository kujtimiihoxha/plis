package module

import (
	"github.com/spf13/afero"
	"os"
)

func ReadFile(path string, fs afero.Fs) (string, error) {
	d, err := afero.ReadFile(fs, path)
	return string(d), err
}
func WriteFile(data string, path string, fs afero.Fs) error {
	return afero.WriteFile(fs, path, []byte(data), os.ModePerm)
}
