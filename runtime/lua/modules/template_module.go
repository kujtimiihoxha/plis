package modules

import (
	"github.com/kujtimiihoxha/plis/api"
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/yuin/gopher-lua"
)

type TemplatesModule struct {
	templatesAPI *api.TemplateAPI
}

func (t *TemplatesModule) copyTemplate(L *lua.LState) int {
	tplName := L.CheckString(1)
	tplDestination := L.CheckString(2)
	tplModel := L.CheckTable(3)
	model := map[string]interface{}{}
	tplModel.ForEach(func(key lua.LValue, value lua.LValue) {
		model[helpers.ToCamelCaseOrUnderscore(helpers.ToGoValue(key).(string))] = helpers.ToGoValue(value)
	})
	err := t.templatesAPI.CopyTemplate(tplName, tplDestination, model)
	if err != nil {
		L.RaiseError("Could not copy template : '%s'", err)
		return 0
	}
	return 0
}
func (t *TemplatesModule) copyTemplateFolder(L *lua.LState) int {
	tplFolder := L.CheckString(1)
	tplDestination := L.ToString(2)
	tplModel := L.ToTable(3)
	excludes := L.ToTable(4)
	exFiles := []string{}
	if excludes != nil {
		for _, v := range helpers.ToGoValue(excludes).([]interface{}) {
			exFiles = append(exFiles, v.(string))
		}
	}
	model := map[string]interface{}{}
	tplModel.ForEach(func(key lua.LValue, value lua.LValue) {
		model[helpers.ToCamelCaseOrUnderscore(helpers.ToGoValue(key).(string))] = helpers.ToGoValue(value)
	})
	err := t.templatesAPI.CopyTemplateFolder(tplFolder, tplDestination, model, exFiles)
	if err != nil {
		L.RaiseError("Could not copy template : '%s'", err)
		return 0
	}
	return 0
}
func (t *TemplatesModule) ModuleLoader() func(L *lua.LState) int {
	return func(L *lua.LState) int {
		mod := L.SetFuncs(L.NewTable(), t.InitializeModule())
		L.Push(mod)
		return 1
	}
}
func (t *TemplatesModule) InitializeModule() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"copyTemplate":       t.copyTemplate,
		"copyTemplateFolder": t.copyTemplateFolder,
	}
}

func NewTemplatesModule(templatesAPI *api.TemplateAPI) *TemplatesModule {
	return &TemplatesModule{
		templatesAPI: templatesAPI,
	}
}
