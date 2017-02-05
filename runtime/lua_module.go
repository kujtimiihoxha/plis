package runtime

import (
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/kujtimiihoxha/plis/module"
	"github.com/spf13/afero"
	"github.com/yuin/gopher-lua"
	"regexp"
	"bytes"
	"encoding/json"
	"github.com/kujtimiihoxha/plis/helpers"
)

func copyTemplate(lr LuaRuntime, L *lua.LState) int {
	tplName := L.ToString(1)
	tplDestination := L.ToString(2)
	tplModel := L.ToTable(3)
	v, err := module.ReadFile(tplName, afero.NewBasePathFs(lr.gFs, "templates"))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	tpl, err := pongo2.FromString(v)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	model := map[string]interface{}{}
	tplModel.ForEach(func(key lua.LValue, value lua.LValue) {
		model[toCamelCase(toGoValue(key).(string))] = toGoValue(value)
	})
	fmt.Println(model)
	out, err := tpl.Execute(pongo2.Context(model))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	err = module.WriteFile(out, tplDestination, fs.GetCurrentFs())
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}
func toJsonFile(L *lua.LState) int {
	destination := L.ToString(1)
	lModel := L.ToTable(2)
	model := toGoValue(lModel)
	err := module.ToJsonFile(destination,model,fs.GetCurrentFs())
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func jsonDecode(L *lua.LState) int {
	str := L.CheckString(1)

	var value interface{}
	err := json.Unmarshal([]byte(str), &value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(helpers.FromJSON(L, value))
	return 1
}

func jsonEncode(L *lua.LState) int {
	value := L.CheckAny(1)

	visited := make(map[*lua.LTable]bool)
	data, err := helpers.ToJSON(value, visited)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(data)))
	return 1
}
func jsonEncodeF(L *lua.LState) int {
	value := L.CheckAny(1)

	visited := make(map[*lua.LTable]bool)
	data, err := helpers.ToJSONFormat(value, visited)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(data)))
	return 1
}
func toGoValue(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[string]interface{})
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(toGoValue(key))
				ret[toCamelCase(keystr)] = toGoValue(value)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, toGoValue(v.RawGetInt(i)))
			}
			return ret
		}
	default:
		return v
	}
}
var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func toCamelCase(src string)(string){
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 { chunks[idx] = bytes.Title(val) }
	}
	return string(bytes.Join(chunks, nil))
}

func readFile(L *lua.LState) int {
	tplName := L.ToString(1)
	v, err := module.ReadFile(tplName, fs.GetCurrentFs())
	L.Push(lua.LString(v))
	if err != nil {
		L.Push(lua.LString(err.Error()))
	}
	L.Push(lua.LNil)
	return 2
}
func ModuleLoader(lr LuaRuntime) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		mod := L.SetFuncs(L.NewTable(), InitializeModule(lr))
		L.Push(mod)
		return 1
	}
}
func InitializeModule(lr LuaRuntime) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"copyTemplate": func(l *lua.LState) int {
			return copyTemplate(lr, l)
		},
		"readFile": readFile,
		"toJsonFile": toJsonFile,
		"jsonDecode": jsonDecode,
		"jsonEncode": jsonEncode,
		"jsonEncodeF": jsonEncodeF,
	}
}
