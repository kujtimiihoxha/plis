package api

import (
	"bytes"
	"github.com/flosch/pongo2"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/ryanuber/go-glob"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type TemplateAPI struct {
	templateFs *FsAPI
	currentFs  *FsAPI
}
type PlisTemplateLoader struct {
	fsAPI *FsAPI
}

func (ptl *PlisTemplateLoader) Abs(base, name string) string {
	return name
}
func (ptl *PlisTemplateLoader) Get(path string) (io.Reader, error) {
	b, err := afero.ReadFile(ptl.fsAPI.fs, path)
	r := bytes.NewReader(b)
	return r, err
}
func (t *TemplateAPI) newTemplateLoader() *PlisTemplateLoader {
	return &PlisTemplateLoader{
		fsAPI: t.templateFs,
	}
}
func (t *TemplateAPI) ReadTemplate(name string, model map[string]interface{}) (string, error) {
	tpSet := pongo2.NewSet("plis_set", t.newTemplateLoader())
	v, err := t.templateFs.ReadFile(name)
	if err != nil {
		return "", err
	}
	if len(model) == 0 {
		return v, nil
	}
	tpl, err := tpSet.FromString(v)
	if err != nil {
		return "", err
	}
	out, err := tpl.Execute(pongo2.Context(model))
	if err != nil {
		return "", err
	}
	return out, nil
}
func (t *TemplateAPI) CopyTemplate(name string, destination string, model map[string]interface{}) error {
	tpSet := pongo2.NewSet("plis_set", t.newTemplateLoader())
	v, err := t.templateFs.ReadFile(name)
	if err != nil {
		return err
	}
	if len(model) == 0 {
		err = t.currentFs.WriteFile(destination, v, false)
		if err != nil {
			return err
		}
		return nil
	}
	tpl, err := tpSet.FromString(v)
	if err != nil {
		return err
	}
	out, err := tpl.Execute(pongo2.Context(model))
	if err != nil {
		return err
	}
	err = t.currentFs.WriteFile(destination, out, false)
	if err != nil {
		return err
	}
	return nil
}

func (t *TemplateAPI) CopyTemplateFolder(folder string, destination string, model map[string]interface{}, excludes []string) error {
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
			v = filepath.ToSlash(v)
			p := filepath.ToSlash(path)
			if glob.Glob(v, p) {
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
