package js

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/robertkrimen/otto"
	"github.com/kujtimiihoxha/plis/config"
	"github.com/spf13/pflag"
	"fmt"
	"github.com/kujtimiihoxha/plis/runtime/js/modules"
	"github.com/kujtimiihoxha/plis/logger"
	"strconv"
	"github.com/kujtimiihoxha/plis/api"
	"github.com/kujtimiihoxha/plis/fs"
)

type Runtime struct {
	vm *otto.Otto
	fs  afero.Fs
	cmd *cobra.Command
	modules map[string]*otto.Object
}

func (js *Runtime) Initialize(cmd *cobra.Command, args map[string]string, c config.GeneratorConfig){
	js.cmd = cmd
	js.vm = otto.New()
	flags,_ := js.vm.Call("new Object",nil)
	getFlags(flags.Object(),js.cmd.Flags())
	a,_ := js.vm.Call("new Object",nil)
	getArguments(a.Object(),c.Args,args)
	js.modules =map[string]*otto.Object{}
	js.modules["plis"] = modules.NewPlisModule(flags.Object(),a.Object(),api.NewPlisAPI(js.cmd)).ModuleLoader(js.vm)
	js.modules["fileSystem"] = modules.NewFileSystemModule(api.NewFsAPI(fs.GetCurrentFs())).ModuleLoader(js.vm)
	js.modules["json"] = modules.NewJSONModule().ModuleLoader(js.vm)
	js.modules["template"] = modules.NewTemplatesModule(
		api.NewTemplatesAPI(
			api.NewFsAPI(afero.NewBasePathFs(js.fs, "templates")),
			api.NewFsAPI(fs.GetCurrentFs()),
		),
	).ModuleLoader(js.vm)
	js.vm.Set("require",js.require)
}
func (js *Runtime) Run() error {
	d, err := afero.ReadFile(js.fs, "run.js")
	if err != nil {
		logger.GetLogger().Error("Could not read run file")
		return err
	}
	if _,err :=js.vm.Run(string(d)); err != nil {
		logger.GetLogger().Fatal(err)
	}
	if v,err:=js.vm.Call("main",nil); err != nil{
		logger.GetLogger().Fatal(err)
	} else {
		if !(v.IsNaN() || v.IsNull() || v.IsUndefined()){
			logger.GetLogger().Fatal(v.String())
		}
	}
	return nil
}
func getArguments(tb *otto.Object, args []config.GeneratorArgs, argsMap map[string]string) {
	for _, v := range args {
		if argsMap[v.Name] == "" && v.Required == false {
			switch v.Type {
			case "string":
				tb.Set(v.Name, "")
			case "int":
				tb.Set(v.Name, 0)
			case "float":
				tb.Set(v.Name, 0.0)
			case "bool":
				tb.Set(v.Name, false)
			}
			tb.Set(v.Name, nil)
			continue
		}
		tb.Set(v.Name, getArgumentByType(v.Type, argsMap[v.Name], v.Name))
	}
}
func getArgumentByType(tp string, arg string, name string) interface {}{
	switch tp {
	case "string":
		return arg
	case "int":
		v, err := strconv.ParseInt(arg, 10, 0)
		if err != nil {
			logger.GetLogger().Fatalf("Argument '%s' must be int", name)
		}
		return  v
	case "float":
		v, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			logger.GetLogger().Fatalf("Argument '%s' must be float", name)
		}
		return v
	case "bool":
		v, err := strconv.ParseBool(arg)
		if err != nil {
			logger.GetLogger().Fatalf("Argument '%s' must be bool", name)
		}
		return v
	}
	return nil
}

func getFlags(tb *otto.Object, flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		fmt.Println(f.Value.Type())
		tb.Set(f.Name, getFlagByType(f.Value.Type(), f.Name, flags))
	})
}
func getFlagByType(tp string, flag string, flagSet *pflag.FlagSet) interface{}{
	switch tp {
	case "string":
		v, err := flagSet.GetString(flag)
		if err != nil {
			return ""
		}
		return v
	case "int":
		v, err := flagSet.GetInt(flag)
		if err != nil {
			return 0
		}
		return v
	case "float64":
		v, err := flagSet.GetFloat64(flag)
		if err != nil {
			return 0.0
		}
		return v
	case "bool":
		v, err := flagSet.GetBool(flag)
		if err != nil {
			return false
		}
		return v
	}
	return nil
}
func NewJsRuntime(fs afero.Fs) *Runtime {
	return &Runtime{
		fs: fs,
	}
}
func (js *Runtime) require(call otto.FunctionCall) otto.Value {
	file := call.Argument(0).String()
	if m:=js.modules[file]; m != nil{
		return m.Value()
	}
	ex,_ := afero.Exists(js.fs,file)
	if !ex {
		ex,_ = afero.Exists(js.fs, file + ".js")
		if ex {
			file = file + ".js"
		}
	}
	data, err := afero.ReadFile(js.fs,file)
	if err != nil {
		return otto.UndefinedValue()
	}
	v, err := call.Otto.Run(string(data))
	if err != nil {
		return otto.UndefinedValue()
	}
	return v.Object().Value()
}