package install

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/aquaproj/aqua/pkg/checksum"
	"github.com/aquaproj/aqua/pkg/config"
	finder "github.com/aquaproj/aqua/pkg/config-finder"
	reader "github.com/aquaproj/aqua/pkg/config-reader"
	"github.com/aquaproj/aqua/pkg/config/aqua"
	registry "github.com/aquaproj/aqua/pkg/install-registry"
	"github.com/aquaproj/aqua/pkg/installpackage"
	"github.com/aquaproj/aqua/pkg/policy"
	"github.com/aquaproj/aqua/pkg/runtime"
	"github.com/aquaproj/aqua/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Controller struct {
	packageInstaller   installpackage.Installer
	rootDir            string
	configFinder       ConfigFinder
	configReader       reader.ConfigReader
	registryInstaller  registry.Installer
	fs                 afero.Fs
	runtime            *runtime.Runtime
	tags               map[string]struct{}
	excludedTags       map[string]struct{}
	policyConfigReader policy.ConfigReader
	policyConfigFinder policy.ConfigFinder
	policyValidator    policy.Validator
	skipLink           bool
	requireChecksum    bool
}

func New(param *config.Param, configFinder ConfigFinder, configReader reader.ConfigReader, registInstaller registry.Installer, pkgInstaller installpackage.Installer, fs afero.Fs, rt *runtime.Runtime, policyConfigReader policy.ConfigReader, policyConfigFinder policy.ConfigFinder, policyValidator policy.Validator) *Controller {
	return &Controller{
		rootDir:            param.RootDir,
		configFinder:       configFinder,
		configReader:       configReader,
		registryInstaller:  registInstaller,
		packageInstaller:   pkgInstaller,
		fs:                 fs,
		runtime:            rt,
		skipLink:           param.SkipLink,
		tags:               param.Tags,
		excludedTags:       param.ExcludedTags,
		policyConfigReader: policyConfigReader,
		policyConfigFinder: policyConfigFinder,
		policyValidator:    policyValidator,
		requireChecksum:    param.RequireChecksum,
	}
}

func (ctrl *Controller) Install(ctx context.Context, logE *logrus.Entry, param *config.Param) error { //nolint:cyclop
	if param.Dest == "" { //nolint:nestif
		rootBin := filepath.Join(ctrl.rootDir, "bin")
		if err := util.MkdirAll(ctrl.fs, rootBin); err != nil {
			return fmt.Errorf("create the directory: %w", err)
		}
		if ctrl.runtime.GOOS == "windows" {
			if err := util.MkdirAll(ctrl.fs, filepath.Join(ctrl.rootDir, "bat")); err != nil {
				return fmt.Errorf("create the directory: %w", err)
			}
		}

		if err := ctrl.packageInstaller.InstallProxy(ctx, logE); err != nil {
			return err //nolint:wrapcheck
		}
	}

	policyFile, err := ctrl.policyConfigFinder.Find("", param.PWD)
	if err != nil {
		return fmt.Errorf("find a policy file: %w", err)
	}
	if policyFile != "" {
		if err := ctrl.policyValidator.Validate(policyFile); err != nil {
			if err := ctrl.policyValidator.Warn(logE, policyFile); err != nil {
				logE.WithError(err).Warn("warn an disallowed policy file")
			}
		} else {
			param.PolicyConfigFilePaths = append(param.PolicyConfigFilePaths, policyFile)
		}
	}

	policyCfgs, err := ctrl.policyConfigReader.Read(param.PolicyConfigFilePaths, param.DisablePolicy)
	if err != nil {
		return fmt.Errorf("read policy files: %w", err)
	}

	for _, cfgFilePath := range ctrl.configFinder.Finds(param.PWD, param.ConfigFilePath) {
		if err := ctrl.install(ctx, logE, cfgFilePath, policyCfgs); err != nil {
			return err
		}
	}

	return ctrl.installAll(ctx, logE, param, policyCfgs)
}

func (ctrl *Controller) installAll(ctx context.Context, logE *logrus.Entry, param *config.Param, policyConfigs []*policy.Config) error {
	if !param.All {
		return nil
	}
	for _, cfgFilePath := range param.GlobalConfigFilePaths {
		if _, err := ctrl.fs.Stat(cfgFilePath); err != nil {
			continue
		}
		if err := ctrl.install(ctx, logE, cfgFilePath, policyConfigs); err != nil {
			return err
		}
	}
	return nil
}

func (ctrl *Controller) install(ctx context.Context, logE *logrus.Entry, cfgFilePath string, policyConfigs []*policy.Config) error {
	cfg := &aqua.Config{}
	if cfgFilePath == "" {
		return finder.ErrConfigFileNotFound
	}
	if err := ctrl.configReader.Read(cfgFilePath, cfg); err != nil {
		return err //nolint:wrapcheck
	}

	var checksums *checksum.Checksums
	if cfg.ChecksumEnabled() {
		checksums = checksum.New()
		checksumFilePath, err := checksum.GetChecksumFilePathFromConfigFilePath(ctrl.fs, cfgFilePath)
		if err != nil {
			return err //nolint:wrapcheck
		}
		if err := checksums.ReadFile(ctrl.fs, checksumFilePath); err != nil {
			return fmt.Errorf("read a checksum JSON: %w", err)
		}
		defer func() {
			if err := checksums.UpdateFile(ctrl.fs, checksumFilePath); err != nil {
				logE.WithError(err).Error("update a checksum file")
			}
		}()
	}

	registryContents, err := ctrl.registryInstaller.InstallRegistries(ctx, logE, cfg, cfgFilePath, checksums)
	if err != nil {
		return err //nolint:wrapcheck
	}

	return ctrl.packageInstaller.InstallPackages(ctx, logE, &installpackage.ParamInstallPackages{ //nolint:wrapcheck
		Config:          cfg,
		Registries:      registryContents,
		ConfigFilePath:  cfgFilePath,
		SkipLink:        ctrl.skipLink,
		Tags:            ctrl.tags,
		ExcludedTags:    ctrl.excludedTags,
		PolicyConfigs:   policyConfigs,
		Checksums:       checksums,
		RequireChecksum: ctrl.requireChecksum,
	})
}
