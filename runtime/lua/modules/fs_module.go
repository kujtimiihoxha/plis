package modules

import (
	"fmt"
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
		"readFile":      fsm.readFile,
		"writeFile":     fsm.writeFile,
		"mkdir":         fsm.mkdir,
		"mkdirAll":      fsm.mkdirAll,
		"fileSeparator": fsm.fileSeparator,
		"exists":        fsm.exists,
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
	err := fsm.fsAPI.WriteFile(path, data)
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
