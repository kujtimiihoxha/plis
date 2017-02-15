package modules

import (
	"fmt"
	"github.com/kujtimiihoxha/plis/api"
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/yuin/gopher-lua"
)

type PlisModule struct {
	plisAPI *api.PlisAPI
	flags   *lua.LTable
	args    *lua.LTable
}

func (p *PlisModule) ModuleLoader() func(L *lua.LState) int {
	return func(L *lua.LState) int {
		mod := L.SetFuncs(L.NewTable(), p.InitializeModule())
		L.SetField(mod, "flags", p.flags)
		L.SetField(mod, "args", p.args)
		L.Push(mod)
		return 1
	}
}
func (p *PlisModule) InitializeModule() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"help":       p.help,
		"runPlisCmd": p.runPlisCmd,
	}
}
func (p *PlisModule) help(L *lua.LState) int {
	p.plisAPI.Help()
	return 0
}
func (p *PlisModule) runPlisCmd(L *lua.LState) int {
	pCmd := L.CheckString(1)
	args := L.CheckAny(2)
	v, ok := helpers.ToGoValue(args).([]interface{})
	if !ok {
		L.Push(lua.LString("The arguments must be an array"))
		return 1
	}
	s := []string{}
	for _, a := range v {
		s = append(s, fmt.Sprint(a))
	}
	if err := p.plisAPI.RunPlisCmd(pCmd, s); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
func NewPlisModule(flags *lua.LTable, args *lua.LTable, api *api.PlisAPI) *PlisModule {
	return &PlisModule{
		plisAPI: api,
		flags:   flags,
		args:    args,
	}
}
