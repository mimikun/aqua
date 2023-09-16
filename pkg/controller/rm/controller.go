package rm

import (
	"github.com/aquaproj/aqua/v2/pkg/config"
	reader "github.com/aquaproj/aqua/v2/pkg/config-reader"
	rgst "github.com/aquaproj/aqua/v2/pkg/install-registry"
	"github.com/aquaproj/aqua/v2/pkg/runtime"
	"github.com/spf13/afero"
)

type Controller struct {
	rootDir           string
	fs                afero.Fs
	runtime           *runtime.Runtime
	configFinder      ConfigFinder
	configReader      reader.ConfigReader
	registryInstaller rgst.Installer
}

func New(param *config.Param, fs afero.Fs, rt *runtime.Runtime, configFinder ConfigFinder, configReader reader.ConfigReader, registryInstaller rgst.Installer) *Controller {
	return &Controller{
		rootDir:           param.RootDir,
		fs:                fs,
		runtime:           rt,
		configFinder:      configFinder,
		configReader:      configReader,
		registryInstaller: registryInstaller,
	}
}

type ConfigFinder interface {
	Find(wd, configFilePath string, globalConfigFilePaths ...string) (string, error)
}
