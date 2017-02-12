package lua

import (
	"fmt"
	"github.com/kujtimiihoxha/plis/api"
	"github.com/kujtimiihoxha/plis/config"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/kujtimiihoxha/plis/runtime"
	"github.com/kujtimiihoxha/plis/runtime/lua/modules"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	glua "github.com/yuin/gopher-lua"
	"strconv"
)

type Runtime struct {
	l   *glua.LState
	fs  afero.Fs
	cmd *cobra.Command
}

func (lr *Runtime) Initialize(cmd *cobra.Command, args map[string]string, c config.GeneratorConfig) runtime.RTime {
	lr.cmd = cmd
	flags := glua.LTable{}
	getFlags(&flags, lr.cmd.Flags())
	argsTb := glua.LTable{}
	getArguments(&argsTb, c.Args, args)
	lr.l = glua.NewState()
	lr.l.PreloadModule("plis", modules.NewPlisModule(&flags, &argsTb, api.NewPlisAPI(lr.cmd)).ModuleLoader())
	lr.l.PreloadModule("fileSystem", modules.NewFileSystemModule(api.NewFsAPI(fs.GetCurrentFs())).ModuleLoader())
	lr.l.PreloadModule("json", modules.NewJSONModule().ModuleLoader())
	lr.l.PreloadModule("templates", modules.NewTemplatesModule(
		api.NewTemplatesAPI(
			api.NewFsAPI(afero.NewBasePathFs(lr.fs, "templates")),
			api.NewFsAPI(fs.GetCurrentFs()),
		),
	).ModuleLoader())
	return lr
}
func (lr *Runtime) Run() error {
	d, err := afero.ReadFile(lr.fs, "run.lua")
	if err != nil {
		logger.GetLogger().Error("Could not read run file")
		return err
	}
	defer lr.l.Close()
	script := fmt.Sprintf(
		"package.path =\"%s\" .. [[%s?.lua]]",
		viper.GetString(fmt.Sprintf("plis.generators.%s.root", lr.cmd.Name())),
		afero.FilePathSeparator,
	)
	script += "\n" + string(d)
	if err := lr.l.DoString(script); err != nil {
		logger.GetLogger().Fatal(err)
	}
	if err := lr.l.CallByParam(glua.P{
		Fn:      lr.l.GetGlobal("main"),
		NRet:    1,
		Protect: true,
	}); err != nil {
		logger.GetLogger().Fatal(err)
	}
	ret := lr.l.Get(-1) // returned value
	lr.l.Pop(1)         // remove received value
	if ret != glua.LNil {
		logger.GetLogger().Fatal(ret)
	}
	return err
}
func NewLuaRuntime(fs afero.Fs) *Runtime {
	return &Runtime{
		fs: fs,
	}
}
func getArguments(tb *glua.LTable, args []config.GeneratorArgs, argsMap map[string]string) {
	for _, v := range args {
		if argsMap[v.Name] == "" && v.Required == false {
			switch v.Type {
			case "string":
				tb.RawSet(glua.LString(v.Name), glua.LString(""))
			case "int":
				tb.RawSet(glua.LString(v.Name), glua.LNumber(0))
			case "float":
				tb.RawSet(glua.LString(v.Name), glua.LNumber(0.0))
			case "bool":
				tb.RawSet(glua.LString(v.Name), glua.LBool(false))
			}
			tb.RawSet(glua.LString(v.Name), glua.LNil)
			continue
		}
		tb.RawSet(glua.LString(v.Name), getArgumentByType(v.Type, argsMap[v.Name], v.Name))
	}
}
func getArgumentByType(tp string, arg string, name string) glua.LValue {
	switch tp {
	case "string":
		return glua.LString(arg)
	case "int":
		v, err := strconv.ParseInt(arg, 10, 0)
		if err != nil {
			logger.GetLogger().Fatalf("Argument '%s' must be int", name)
		}
		return glua.LNumber(v)
	case "float":
		v, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			logger.GetLogger().Fatalf("Argument '%s' must be float", name)
		}
		return glua.LNumber(v)
	case "bool":
		v, err := strconv.ParseBool(arg)
		if err != nil {
			logger.GetLogger().Fatalf("Argument '%s' must be bool", name)
		}
		return glua.LBool(v)
	}
	return glua.LNil
}
func getFlags(tb *glua.LTable, flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		tb.RawSet(glua.LString(f.Name), getFlagByType(f.Value.Type(), f.Name, flags))
	})
}
func getFlagByType(tp string, flag string, flagSet *pflag.FlagSet) glua.LValue {
	switch tp {
	case "string":
		v, err := flagSet.GetString(flag)
		if err != nil {
			return glua.LString("")
		}
		return glua.LString(v)
	case "int":
		v, err := flagSet.GetInt(flag)
		if err != nil {
			return glua.LNumber(0)
		}
		return glua.LNumber(v)
	case "float":
		v, err := flagSet.GetInt(flag)
		if err != nil {
			return glua.LNumber(0)
		}
		return glua.LNumber(v)
	case "bool":
		v, err := flagSet.GetBool(flag)
		if err != nil {
			return glua.LBool(false)
		}
		return glua.LBool(v)
	}
	return glua.LNil
}
