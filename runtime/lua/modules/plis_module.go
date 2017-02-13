package modules

import (
	"fmt"
	"github.com/kujtimiihoxha/plis/api"
	"github.com/kujtimiihoxha/plis/cmd"
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
	c := L.CheckString(1)
	args := L.CheckAny(2)
	v, ok := helpers.ToGoValue(args).([]interface{})
	if !ok {
		L.RaiseError("The arguments must be an array")
		return 0
	}
	s := []string{
		c,
	}
	for _, a := range v {
		s = append(s, fmt.Sprint(a))
	}
	cmd.RootCmd.SetArgs(s)
	if err := cmd.RootCmd.Execute(); err != nil {
		L.RaiseError(err.Error())
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
