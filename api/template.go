package api

import (
	"github.com/flosch/pongo2"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/ryanuber/go-glob"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
)

type TemplateAPI struct {
	templateFs *FsAPI
	currentFs  *FsAPI
}

func (t *TemplateAPI) CopyTemplate(name string, destination string, model map[string]interface{}) error {
	name = strings.Replace(name,"\\",t.templateFs.FileSeparator(),-1)
	destination = strings.Replace(destination,"\\",t.templateFs.FileSeparator(),-1)
	v, err := t.templateFs.ReadFile(name)
	if err != nil {
		return err
	}
	if len(model) == 0 {
		err = t.currentFs.WriteFile(destination, v)
		if err != nil {
			return err
		}
		return nil
	}
	tpl, err := pongo2.FromString(v)
	if err != nil {
		return err
	}
	out, err := tpl.Execute(pongo2.Context(model))
	if err != nil {
		return err
	}
	err = t.currentFs.WriteFile(destination, out)
	if err != nil {
		return err
	}
	return nil
}

func (t *TemplateAPI) CopyTemplateFolder(folder string, destination string, model map[string]interface{}, excludes []string) error {
	folder = strings.Replace(folder,"\\",t.templateFs.FileSeparator(),-1)
	destination = strings.Replace(destination,"\\",t.templateFs.FileSeparator(),-1)
	if destination != "" {
		b, err := afero.Exists(t.currentFs.fs, destination)
		if err != nil {
			return err
		}
		if !b {
			err = t.currentFs.MkdirAll(destination)
			if err != nil {
				return err
			}
		}
		t.currentFs.fs = afero.NewBasePathFs(fs.GetCurrentFs(), destination)
	}
	err := afero.Walk(t.templateFs.fs, "", func(path string, info os.FileInfo, err error) error {
		for _, v := range excludes {
			if glob.Glob(v, path) {
				return nil
			}
		}
		if !info.IsDir() {
			destName := strings.TrimSuffix(path, ".tpl")
			err = t.currentFs.MkdirAll(filepath.Dir(destName))
			if err != nil {
				return err
			}
			err = t.CopyTemplate(path, destName, model)
			if err != nil {
				return err
			}
		} else {
			err = t.currentFs.MkdirAll(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func NewTemplatesAPI(templateFs *FsAPI, currentFs *FsAPI) *TemplateAPI {
	return &TemplateAPI{
		templateFs: templateFs,
		currentFs:  currentFs,
	}
}
