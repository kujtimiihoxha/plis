package helpers

import (
	"errors"
	"fmt"
	"github.com/yuin/gopher-lua"
	"strconv"
)

var (
	errFunction = errors.New("cannot encode function to JSON")
	errChannel  = errors.New("cannot encode channel to JSON")
	errState    = errors.New("cannot encode state to JSON")
	errUserData = errors.New("cannot encode userdata to JSON")
	errNested   = errors.New("cannot encode recursively nested tables to JSON")
)

type jsonValue struct {
	lua.LValue
	marshallFunc func(v interface{}) ([]byte, error)
	visited      map[*lua.LTable]bool
}

func (j jsonValue) MarshalJSON() ([]byte, error) {
	return ToJson(j.LValue, j.visited, j.marshallFunc)
}
func ToGoValue(lv lua.LValue) interface{} {
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
				keystr := fmt.Sprint(ToGoValue(key))
				ret[ToCamelCaseOrUnderscore(keystr)] = ToGoValue(value)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, ToGoValue(v.RawGetInt(i)))
			}
			return ret
		}
	default:
		return v
	}
}
func ToJson(value lua.LValue, visited map[*lua.LTable]bool, marshallFunc func(v interface{}) ([]byte, error)) (data []byte, err error) {
	switch converted := value.(type) {
	case lua.LBool:
		data, err = marshallFunc(converted)
	case lua.LChannel:
		err = errChannel
	case lua.LNumber:
		data, err = marshallFunc(converted)
	case *lua.LFunction:
		err = errFunction
	case *lua.LNilType:
		data, err = marshallFunc(converted)
	case *lua.LState:
		err = errState
	case lua.LString:
		data, err = marshallFunc(converted)
	case *lua.LTable:
		var arr []jsonValue
		var obj map[string]jsonValue

		if visited[converted] {
			panic(errNested)
		}
		visited[converted] = true

		converted.ForEach(func(k lua.LValue, v lua.LValue) {
			i, numberKey := k.(lua.LNumber)
			if numberKey && obj == nil {
				index := int(i) - 1
				if index != len(arr) {
					// map out of order; convert to map
					obj = make(map[string]jsonValue)
					for i, value := range arr {
						obj[strconv.Itoa(i+1)] = value
					}
					obj[strconv.Itoa(index+1)] = jsonValue{v, marshallFunc, visited}
					return
				}
				arr = append(arr, jsonValue{v, marshallFunc, visited})
				return
			}
			if obj == nil {
				obj = make(map[string]jsonValue)
				for i, value := range arr {
					obj[strconv.Itoa(i+1)] = value
				}
			}
			obj[k.String()] = jsonValue{v, marshallFunc, visited}
		})
		if obj != nil {
			data, err = marshallFunc(obj)
		} else {
			data, err = marshallFunc(arr)
		}
	case *lua.LUserData:
		// TODO: call metatable __tostring?
		err = errUserData
	}
	return
}
func FromJson(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(FromJson(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), FromJson(L, item))
		}
		return tbl
	}
	return lua.LNil
}
