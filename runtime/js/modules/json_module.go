package modules

import (
	"github.com/robertkrimen/otto"
)

type JSONModule struct{}

func (j *JSONModule) ModuleLoader(vm *otto.Otto) *otto.Object {
	obj,_ := vm.Call("new Object",nil)
	v := obj.Object()
	v.Set("decode",j.decode)
	v.Set("encode",j.encode)
	v.Set("encodeF",j.encodeF)
	return v
}
func (j *JSONModule) decode(call otto.FunctionCall) otto.Value  {
	obj,_ := call.Otto.Call("JSON.parse",nil,call.Argument(0))
	return obj
}
func (j *JSONModule) encode(call otto.FunctionCall) otto.Value  {
	obj,_ := call.Otto.Call("JSON.stringify",nil,call.Argument(0))
	return obj
}
func (j *JSONModule) encodeF(call otto.FunctionCall) otto.Value  {
	obj,_ := call.Otto.Call("JSON.stringify",nil,call.Argument(0),nil,2)
	return obj
}
func NewJSONModule() *JSONModule {
	return &JSONModule{}
}