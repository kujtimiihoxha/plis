package module

import (
	"github.com/spf13/afero"
	"os"
	"encoding/json"
)

func ReadFile(path string, fs afero.Fs) (string, error) {
	d, err := afero.ReadFile(fs, path)
	return string(d), err
}
func WriteFile(data string, path string, fs afero.Fs) error {
	return afero.WriteFile(fs, path, []byte(data), os.ModePerm)
}
func ToJsonFile(path string, m interface{}, fs afero.Fs)  error{
	d, err :=json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}
	return afero.WriteFile(fs, path, d, os.ModePerm)
}