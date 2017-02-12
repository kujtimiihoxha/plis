package modules

import (
	"encoding/json"
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/yuin/gopher-lua"
)

type JSONModule struct{}

func (j *JSONModule) ModuleLoader() func(L *lua.LState) int {
	return func(L *lua.LState) int {
		mod := L.SetFuncs(L.NewTable(), j.InitializeModule())
		L.Push(mod)
		return 1
	}
}
func (j *JSONModule) InitializeModule() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"encode":  j.encode,
		"encodeF": j.encodeF,
		"decode":  j.decode,
	}
}

func NewJSONModule() *JSONModule {
	return &JSONModule{}
}

func (j *JSONModule) decode(L *lua.LState) int {
	str := L.CheckString(1)
	var value interface{}
	err := json.Unmarshal([]byte(str), &value)
	if err != nil {
		L.Push(lua.LNil)
		L.RaiseError("Could not decode json : '%s'", err)
		return 1
	}
	L.Push(helpers.FromJSON(L, value))
	return 1
}
func (j *JSONModule) encode(L *lua.LState) int {
	value := L.CheckAny(1)
	visited := make(map[*lua.LTable]bool)
	data, err := helpers.ToJSON(value, visited, json.Marshal)
	if err != nil {
		L.Push(lua.LNil)
		L.RaiseError("Could not encode json : '%s'", err)
		return 1
	}
	L.Push(lua.LString(string(data)))
	return 1
}
func (j *JSONModule) encodeF(L *lua.LState) int {
	value := L.CheckAny(1)
	visited := make(map[*lua.LTable]bool)
	data, err := helpers.ToJSON(value, visited, marshalFormat)
	if err != nil {
		L.Push(lua.LNil)
		L.RaiseError("Could not encode json : '%s'", err)
		return 1
	}
	L.Push(lua.LString(string(data)))
	return 1
}
func marshalFormat(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "    ")
}
