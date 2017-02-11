package modules

import (
	"github.com/kujtimiihoxha/plis/api"
	"github.com/yuin/gopher-lua"
)

type PlisModule struct {
	plisApi *api.PlisApi
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
		"help": p.help,
	}
}
func (p *PlisModule) help(L *lua.LState) int {
	p.plisApi.Help()
	return 1
}
func NewPlisModule(flags *lua.LTable, args *lua.LTable, api *api.PlisApi) *PlisModule {
	return &PlisModule{
		plisApi: api,
		flags:   flags,
		args:    args,
	}
}
