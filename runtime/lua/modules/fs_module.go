package modules

import (
	"fmt"
	"github.com/kujtimiihoxha/plis/api"
	"github.com/yuin/gopher-lua"
	"os"
)

type FileSystemModule struct {
	fsAPI *api.FsAPI
}

func (fsm *FileSystemModule) ModuleLoader() func(L *lua.LState) int {
	return func(L *lua.LState) int {
		mod := L.SetFuncs(L.NewTable(), fsm.InitializeModule())
		L.Push(mod)
		return 1
	}
}
func (fsm *FileSystemModule) InitializeModule() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"readFile":      fsm.readFile,
		"writeFile":     fsm.writeFile,
		"mkdir":         fsm.mkdir,
		"mkdirAll":      fsm.mkdirAll,
		"fileSeparator": fsm.fileSeparator,
		"exists":        fsm.exists,
		"walk":          fsm.walk,
	}
}

func NewFileSystemModule(fsAPI *api.FsAPI) *FileSystemModule {
	return &FileSystemModule{
		fsAPI: fsAPI,
	}
}

func (fsm *FileSystemModule) readFile(L *lua.LState) int {
	fName := L.CheckString(1)
	v, err := fsm.fsAPI.ReadFile(fName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Could not read file : '%s'", err)))
		return 2
	}
	L.Push(lua.LString(v))
	return 1
}
func (fsm *FileSystemModule) exists(L *lua.LState) int {
	fName := L.CheckString(1)
	v, err := fsm.fsAPI.Exists(fName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Could not check if file exists : '%s'", err)))
		return 2
	}
	L.Push(lua.LBool(v))
	return 1
}
func (fsm *FileSystemModule) fileSeparator(L *lua.LState) int {
	v := fsm.fsAPI.FilePathSeparator()
	L.Push(lua.LString(v))
	return 1
}
func (fsm *FileSystemModule) writeFile(L *lua.LState) int {
	path := L.CheckString(1)
	data := L.CheckString(2)
	b := L.ToBool(3)
	err := fsm.fsAPI.WriteFile(path, data,b)
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("Could not write file : '%s'", err)))
		return 1
	}
	return 0
}
func (fsm *FileSystemModule) mkdir(L *lua.LState) int {
	path := L.CheckString(1)
	err := fsm.fsAPI.Mkdir(path)
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("Could not create directory file : '%s'", err)))
		return 1
	}
	return 0
}
func (fsm *FileSystemModule) mkdirAll(L *lua.LState) int {
	path := L.CheckString(1)
	err := fsm.fsAPI.MkdirAll(path)
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("Could not create directories file : '%s'", err)))
		return 1
	}
	return 0
}

func (fsm *FileSystemModule) walk(L *lua.LState) int {
	root := L.CheckString(1)
	fc := L.CheckFunction(2)
	inf := L.NewTable()
	fsm.fsAPI.Walk(root, func(path string, info os.FileInfo, err error) error {
		inf.RawSet(lua.LString("isDir"), lua.LBool(info.IsDir()))
		inf.RawSet(lua.LString("name"), lua.LString(info.Name()))
		inf.RawSet(lua.LString("size"), lua.LNumber(info.Size()))
		e := ""
		if err != nil {
			e = err.Error()
		}
		err = L.CallByParam(lua.P{
			Fn:      fc,
			NRet:    0,
			Protect: true,
		}, lua.LString(path), inf, lua.LString(e))
		return nil
	})
	return 0
}
