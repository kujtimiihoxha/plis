package modules

import (
	"github.com/kujtimiihoxha/plis/api"
	"github.com/yuin/gopher-lua"
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
		"readFile":  fsm.readFile,
		"writeFile": fsm.writeFile,
		"mkdir":     fsm.mkdir,
		"mkdirAll":  fsm.mkdirAll,
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
	err := fsm.fsAPI.WriteFile(path, data)
	if err != nil {
		L.RaiseError("Could not write file : '%s'", err)
		return 0
	}
	return 0
}
func (fsm *FileSystemModule) mkdir(L *lua.LState) int {
	path := L.CheckString(1)
	err := fsm.fsAPI.Mkdir(path)
	if err != nil {
		L.RaiseError("Could not create directory file : '%s'", err)
		return 0
	}
	return 0
}
func (fsm *FileSystemModule) mkdirAll(L *lua.LState) int {
	path := L.CheckString(1)
	err := fsm.fsAPI.MkdirAll(path)
	if err != nil {
		L.RaiseError("Could not create directories file : '%s'", err)
		return 0
	}
	return 0
}
