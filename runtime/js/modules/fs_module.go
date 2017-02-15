package modules

import (
	"github.com/robertkrimen/otto"
	"github.com/kujtimiihoxha/plis/api"
	"github.com/kujtimiihoxha/plis/logger"
)

type FileSystemModule struct {
	fsAPI *api.FsAPI
}

func (fsm *FileSystemModule) ModuleLoader(vm *otto.Otto) *otto.Object {
	obj,_ := vm.Call("new Object",nil)
	v := obj.Object()
	v.Set("readFile",fsm.readFile)
	v.Set("exists",fsm.exists)
	v.Set("fileSeparator",fsm.fileSeparator)
	v.Set("writeFile",fsm.writeFile)
	v.Set("mkdir",fsm.mkdir)
	v.Set("mkdirAll",fsm.mkdirAll)
	return v
}

func NewFileSystemModule(fsAPI *api.FsAPI) *FileSystemModule {
	return &FileSystemModule{
		fsAPI: fsAPI,
	}
}
func (fsm *FileSystemModule) readFile(call otto.FunctionCall) otto.Value  {
	fName := call.Argument(0).String()
	v, err := fsm.fsAPI.ReadFile(fName)
	if err != nil {
		logger.GetLogger().Errorf("Could not check if file exists : '%s'", err)
		return otto.UndefinedValue()
	}
	obj,_ :=otto.ToValue(v)
	return obj
}
func (fsm *FileSystemModule) writeFile(call otto.FunctionCall) otto.Value {
	path := call.Argument(0).String()
	data := call.Argument(1).String()
	err := fsm.fsAPI.WriteFile(path, data)
	if err != nil {
		logger.GetLogger().Errorf("Could not write file : '%s'", err)
		return otto.UndefinedValue()
	}
	return otto.TrueValue()
}
func (fsm *FileSystemModule) fileSeparator(call otto.FunctionCall) otto.Value{
	v := fsm.fsAPI.FilePathSeparator()
	obj,_ :=otto.ToValue(v)
	return obj
}
func (fsm *FileSystemModule) exists(call otto.FunctionCall) otto.Value {
	fName := call.Argument(0).String()
	v, err := fsm.fsAPI.Exists(fName)
	if err != nil {
		logger.GetLogger().Errorf("Could not check if file exists : '%s'", err)
		return otto.UndefinedValue()
	}
	obj,_ :=otto.ToValue(v)
	return obj
}

func (fsm *FileSystemModule) mkdir(call otto.FunctionCall) otto.Value{
	path := call.Argument(0).String()
	err := fsm.fsAPI.Mkdir(path)
	if err != nil {
		logger.GetLogger().Errorf("Could not create directory file : '%s'", err)
		return otto.UndefinedValue()
	}
	return otto.TrueValue()
}
func (fsm *FileSystemModule) mkdirAll(call otto.FunctionCall) otto.Value{
	path := call.Argument(0).String()
	err := fsm.fsAPI.MkdirAll(path)
	if err != nil {
		logger.GetLogger().Errorf("Could not create directories file : '%s'", err)
		return otto.UndefinedValue()
	}
	return otto.TrueValue()
}
