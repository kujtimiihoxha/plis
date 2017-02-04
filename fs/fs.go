package fs

import (
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var plisRootFs afero.Fs
var currentFs afero.Fs
var generatorTestDir afero.Fs

func initialize() {
	if viper.GetString("plis.dir.root") == "" {
		logger.GetLogger().Fatal("Plis root configuration not set.")
	}
	if !viper.GetBool("plis.testing") {
		plisRootFs = afero.NewBasePathFs(afero.NewOsFs(), viper.GetString("plis.dir.root"))
		currentFs = afero.NewOsFs()
	} else {
		plisRootFs = afero.NewMemMapFs()
		currentFs = afero.NewMemMapFs()
	}
}

func GetPlisRootFs() afero.Fs {
	if plisRootFs == nil {
		initialize()
	}
	return plisRootFs
}

func GetCurrentFs() afero.Fs {
	if currentFs == nil {
		initialize()
	}
	return currentFs
}

func SetGeneratorTestFs(fs afero.Fs) {
	generatorTestDir = fs
}

func GetGeneratorTestFs() afero.Fs {
	return generatorTestDir
}
