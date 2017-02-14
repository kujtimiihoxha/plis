package modules

import (
	"github.com/robertkrimen/otto"
	"github.com/kujtimiihoxha/plis/api"
)

type PlisModule struct {
	plisAPI *api.PlisAPI
	flags   *otto.Object
	args    *otto.Object
}

func (p *PlisModule) ModuleLoader(vm *otto.Otto) *otto.Object {
	obj,_ := vm.Call("new Object",nil)
	v := obj.Object()
	v.Set("help",p.help)
	v.Set("runPlisCmd",p.runPlisCmd)
	v.Set("flags",p.flags)
	v.Set("args",p.args)
	return v
}
func (p *PlisModule) help(call otto.FunctionCall) otto.Value  {
	p.plisAPI.Help()
	return otto.Value{}
}
func (p *PlisModule) runPlisCmd(call otto.FunctionCall) otto.Value  {
	pCmd := call.Argument(0).String()
	args := call.Argument(1)
	if !args.IsObject() {
		v,_:= otto.ToValue("The arguments must be an array")
		return v
	}
	s := []string{}
	for _, a := range args.Object().Keys() {
		vl,_:= args.Object().Get(a)
		s = append(s, vl.String())
	}
	if err := p.plisAPI.RunPlisCmd(pCmd,s); err != nil {
		v,_:= otto.ToValue(err.Error())
		return v
	}
	return otto.Value{}
}

func NewPlisModule(flags *otto.Object, args *otto.Object, api *api.PlisAPI) *PlisModule {
	return &PlisModule{
		plisAPI: api,
		flags:   flags,
		args:    args,
	}
}