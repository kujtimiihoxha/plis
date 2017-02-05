package runtime

import (
	"github.com/kujtimiihoxha/plis/config"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/spf13/afero"
	"github.com/yuin/gopher-lua"
)

type LuaRuntime struct {
	l   *lua.LState
	gFs afero.Fs
}

func (lr LuaRuntime) Initialize(c config.GeneratorConfig, gFs afero.Fs) RTime {
	lr.gFs = gFs
	lr.l = lua.NewState()
	lr.l.PreloadModule("plis", ModuleLoader(lr))
	return lr
}
func (lr LuaRuntime) Run() error {
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
