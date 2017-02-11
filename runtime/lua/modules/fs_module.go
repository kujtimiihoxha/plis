package modules

import (
	"github.com/kujtimiihoxha/plis/api"
	"github.com/yuin/gopher-lua"
)

type FileSystemModule struct {
	fsApi *api.FsApi
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
		"readFile":  fsm.readFile,
		"writeFile": fsm.writeFile,
		"mkdir":     fsm.mkdir,
		"mkdirAll":  fsm.mkdirAll,
	}
}

func NewFileSystemModule(fsApi *api.FsApi) *FileSystemModule {
	return &FileSystemModule{
		fsApi: fsApi,
	}
}

func (fsm *FileSystemModule) readFile(L *lua.LState) int {
	fName := L.CheckString(1)
	v, err := fsm.fsApi.ReadFile(fName)
	if err != nil {
		L.RaiseError("Could not read file : '%s'", err)
		L.Push(lua.LNil)
		return 1
	}
	L.Push(lua.LString(v))
	return 1
}
func (fsm *FileSystemModule) writeFile(L *lua.LState) int {
	path := L.CheckString(1)
	data := L.CheckString(2)
	err := fsm.fsApi.WriteFile(path, data)
	if err != nil {
		L.RaiseError("Could not write file : '%s'", err)
		return 0
	}
	return 0
}
func (fsm *FileSystemModule) mkdir(L *lua.LState) int {
	path := L.CheckString(1)
	err := fsm.fsApi.Mkdir(path)
	if err != nil {
		L.RaiseError("Could not create directory file : '%s'", err)
		return 0
	}
	return 0
}
func (fsm *FileSystemModule) mkdirAll(L *lua.LState) int {
	path := L.CheckString(1)
	err := fsm.fsApi.MkdirAll(path)
	if err != nil {
		L.RaiseError("Could not create directories file : '%s'", err)
		return 0
	}
	return 0
}
