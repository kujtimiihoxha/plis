package runtime

import (
	"github.com/kujtimiihoxha/plis/config"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/yuin/gopher-lua"
	"strconv"
)

type LuaRuntime struct {
	l   *lua.LState
	gFs afero.Fs
}

func (lr LuaRuntime) Initialize(cmd *cobra.Command, args map[string]string, c config.GeneratorConfig, gFs afero.Fs) RTime {
	lr.gFs = gFs
	_cmd := cmd
	flags := lua.LTable{}
	_cmd.Flags().VisitAll(func(f *pflag.Flag) {
		flags.RawSet(lua.LString(f.Name), getFlagByType(f.Value.Type(), f.Name, _cmd.Flags()))
	})
	argsTb := lua.LTable{}
	for _,v:= range c.Args{
		if  args[v.Name] == "" && v.Required == false {
			switch v.Type {
			case "string":
				argsTb.RawSet(lua.LString(v.Name), lua.LString(""))
			case "int":
				argsTb.RawSet(lua.LString(v.Name), lua.LNumber(0))
			case "float":
				argsTb.RawSet(lua.LString(v.Name), lua.LNumber(0.0))
			case "bool":
				argsTb.RawSet(lua.LString(v.Name), lua.LBool(false))
			}
			argsTb.RawSet(lua.LString(v.Name), lua.LNil)
		}
		argsTb.RawSet(lua.LString(v.Name), getArgumentByType(v.Type, args[v.Name], v.Name))
	}
	lr.l = lua.NewState()
	lr.l.PreloadModule("plis", ModuleLoader(lr, &flags, &argsTb))
	return lr
}
func (lr LuaRuntime) Run() error {
	println("TEST")
	d, err := afero.ReadFile(lr.gFs, "run.lua")
	if err != nil {
		logger.GetLogger().Error("Could not read run file")
		return err
	}
	defer lr.l.Close()
	if err := lr.l.DoString(string(d)); err != nil {
		logger.GetLogger().Error(err)
	}
	return err
}
func getArgumentByType(tp string, arg string, name string) lua.LValue {
	switch tp {
	case "string":
		return lua.LString(arg)
	case "int":
		v, err := strconv.ParseInt(arg,10,0)
		if err != nil {
			logger.GetLogger().Fatalf("Argument '%s' must be int", name)
		}
		return lua.LNumber(v)
	case "float":
		v, err := strconv.ParseFloat(arg,64)
		if err != nil {
			logger.GetLogger().Fatalf("Argument '%s' must be float", name)
		}
		return lua.LNumber(v)
	case "bool":
		v, err := strconv.ParseBool(arg)
		if err != nil {
			logger.GetLogger().Fatalf("Argument '%s' must be bool", name)
		}
		return lua.LBool(v)
	}
	return lua.LNil
}
func getFlagByType(tp string, flag string, flagSet *pflag.FlagSet) lua.LValue {
	switch tp {
	case "string":
		v, err := flagSet.GetString(flag)
		if err != nil {
			return lua.LString("")
		}
		return lua.LString(v)
	case "int":
		v, err := flagSet.GetInt(flag)
		if err != nil {
			return lua.LNumber(0)
		}
		return lua.LNumber(v)
	case "float":
		v, err := flagSet.GetInt(flag)
		if err != nil {
			return lua.LNumber(0)
		}
		return lua.LNumber(v)
	case "bool":
		v, err := flagSet.GetBool(flag)
		if err != nil {
			return lua.LBool(false)
		}
		return lua.LBool(v)
	}
	return lua.LNil
}
