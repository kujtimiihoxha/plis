package modules

import (
	"github.com/kujtimiihoxha/plis/api"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/robertkrimen/otto"
	"os"
)

type FileSystemModule struct {
	fsAPI *api.FsAPI
}

func (fsm *FileSystemModule) ModuleLoader(vm *otto.Otto) *otto.Object {
	obj, _ := vm.Call("new Object", nil)
	v := obj.Object()
	v.Set("readFile", fsm.readFile)
	v.Set("exists", fsm.exists)
	v.Set("fileSeparator", fsm.fileSeparator)
	v.Set("writeFile", fsm.writeFile)
	v.Set("mkdir", fsm.mkdir)
	v.Set("walk", fsm.walk)
	v.Set("mkdirAll", fsm.mkdirAll)
	return v
}

func NewFileSystemModule(fsAPI *api.FsAPI) *FileSystemModule {
	return &FileSystemModule{
		fsAPI: fsAPI,
	}
}
func (fsm *FileSystemModule) readFile(call otto.FunctionCall) otto.Value {
	fName := call.Argument(0).String()
	v, err := fsm.fsAPI.ReadFile(fName)
	if err != nil {
		logger.GetLogger().Errorf("Could not check if file exists : '%s'", err)
		return otto.UndefinedValue()
	}
	obj, _ := otto.ToValue(v)
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
func (fsm *FileSystemModule) fileSeparator(call otto.FunctionCall) otto.Value {
	v := fsm.fsAPI.FilePathSeparator()
	obj, _ := otto.ToValue(v)
	return obj
}
func (fsm *FileSystemModule) exists(call otto.FunctionCall) otto.Value {
	fName := call.Argument(0).String()
	v, err := fsm.fsAPI.Exists(fName)
	if err != nil {
		logger.GetLogger().Errorf("Could not check if file exists : '%s'", err)
		return otto.UndefinedValue()
	}
	obj, _ := otto.ToValue(v)
	return obj
}

func (fsm *FileSystemModule) mkdir(call otto.FunctionCall) otto.Value {
	path := call.Argument(0).String()
	err := fsm.fsAPI.Mkdir(path)
	if err != nil {
		logger.GetLogger().Errorf("Could not create directory file : '%s'", err)
		return otto.UndefinedValue()
	}
	return otto.TrueValue()
}
func (fsm *FileSystemModule) mkdirAll(call otto.FunctionCall) otto.Value {
	path := call.Argument(0).String()
	err := fsm.fsAPI.MkdirAll(path)
	if err != nil {
		logger.GetLogger().Errorf("Could not create directories file : '%s'", err)
		return otto.UndefinedValue()
	}
	return otto.TrueValue()
}

func (fsm *FileSystemModule) walk(call otto.FunctionCall) otto.Value{
	root := call.Argument(0).String()
	fc := call.Argument(1)
	if !fc.IsFunction() {
		logger.GetLogger().Errorf("Walk needs a function to call")
		return otto.FalseValue()
	}
	fsm.fsAPI.Walk(root, func(path string, info os.FileInfo, err error) error {
		inf := map[string]interface{}{}
		inf["isDir"]=info.IsDir()
		inf["name"]=info.Name()
		inf["size"]=info.Size()
		v,_:= call.Otto.ToValue(inf)
		e := ""
		if err != nil {
			e = err.Error()
		}
		_,err=fc.Call(fc.Object().Value(),path,v,e)
		return err
	})
	return otto.NullValue()
}