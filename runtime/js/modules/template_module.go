package modules

import (
	"github.com/kujtimiihoxha/plis/api"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/robertkrimen/otto"
)

type TemplatesModule struct {
	templatesAPI *api.TemplateAPI
}

func (t *TemplatesModule) ModuleLoader(vm *otto.Otto) *otto.Object {
	obj, _ := vm.Call("new Object", nil)
	v := obj.Object()
	v.Set("copyTemplate", t.copyTemplate)
	v.Set("copyTemplateFolder", t.copyTemplateFolder)
	v.Set("readTemplate", t.readTemplate)
	return v
}
func (t *TemplatesModule) readTemplate(call otto.FunctionCall) otto.Value {
	tplName := call.Argument(0).String()
	tplModel, _ := call.Argument(2).Export()
	if tplModel == nil {
		logger.GetLogger().Error("You must provide a model")
		return otto.FalseValue()
	}
	model, ok := tplModel.(map[string]interface{})
	if !ok {
		logger.GetLogger().Error("The template model must be an object")
		return otto.FalseValue()
	}
	v, err := t.templatesAPI.ReadTemplate(tplName, model)
	if err != nil {
		logger.GetLogger().Errorf("Error while copying template :%s", err.Error())
		return otto.FalseValue()
	}
	vl, _ := otto.ToValue(v)
	return vl
}
func (t *TemplatesModule) copyTemplate(call otto.FunctionCall) otto.Value {
	tplName := call.Argument(0).String()
	tplDestination := call.Argument(1).String()
	tplModel, _ := call.Argument(2).Export()
	if tplModel == nil {
		logger.GetLogger().Error("You must provide a model")
		return otto.FalseValue()
	}
	model, ok := tplModel.(map[string]interface{})
	if !ok {
		logger.GetLogger().Error("The template model must be an object")
		return otto.FalseValue()
	}
	if tplDestination == "" {
		tplDestination = tplName
	}
	err := t.templatesAPI.CopyTemplate(tplName, tplDestination, model)
	if err != nil {
		logger.GetLogger().Errorf("Error while copying template :%s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}
func (t *TemplatesModule) copyTemplateFolder(call otto.FunctionCall) otto.Value {
	tplFolder := call.Argument(0).String()
	tplDestination := call.Argument(1).String()
	tplModel, _ := call.Argument(2).Export()
	if tplModel == nil {
		logger.GetLogger().Error("You must provide a model")
		return otto.FalseValue()
	}
	model, ok := tplModel.(map[string]interface{})
	if !ok {
		logger.GetLogger().Error("The template model must be an object")
		return otto.FalseValue()
	}
	excludes, _ := call.Argument(3).Export()
	exFiles := []string{}
	if excludes != nil {
		exFiles, ok = excludes.([]string)
		if !ok {
			logger.GetLogger().Error("The exludes object must be an array of strings")
			return otto.FalseValue()
		}
	}
	err := t.templatesAPI.CopyTemplateFolder(tplFolder, tplDestination, model, exFiles)
	if err != nil {
		logger.GetLogger().Errorf("Error while copying template folder :%s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}
func NewTemplatesModule(templatesAPI *api.TemplateAPI) *TemplatesModule {
	return &TemplatesModule{
		templatesAPI: templatesAPI,
	}
}
