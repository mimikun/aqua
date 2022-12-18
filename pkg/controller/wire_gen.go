// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package controller

import (
	"context"
	"github.com/aquaproj/aqua/pkg/checksum"
	"github.com/aquaproj/aqua/pkg/config"
	"github.com/aquaproj/aqua/pkg/config-finder"
	"github.com/aquaproj/aqua/pkg/config-reader"
	"github.com/aquaproj/aqua/pkg/controller/cp"
	exec2 "github.com/aquaproj/aqua/pkg/controller/exec"
	"github.com/aquaproj/aqua/pkg/controller/generate"
	"github.com/aquaproj/aqua/pkg/controller/generate-registry"
	"github.com/aquaproj/aqua/pkg/controller/initcmd"
	"github.com/aquaproj/aqua/pkg/controller/initpolicy"
	"github.com/aquaproj/aqua/pkg/controller/install"
	"github.com/aquaproj/aqua/pkg/controller/list"
	"github.com/aquaproj/aqua/pkg/controller/updateaqua"
	"github.com/aquaproj/aqua/pkg/controller/updatechecksum"
	"github.com/aquaproj/aqua/pkg/controller/which"
	"github.com/aquaproj/aqua/pkg/cosign"
	"github.com/aquaproj/aqua/pkg/download"
	"github.com/aquaproj/aqua/pkg/exec"
	"github.com/aquaproj/aqua/pkg/github"
	"github.com/aquaproj/aqua/pkg/install-registry"
	"github.com/aquaproj/aqua/pkg/installpackage"
	"github.com/aquaproj/aqua/pkg/link"
	"github.com/aquaproj/aqua/pkg/policy"
	"github.com/aquaproj/aqua/pkg/runtime"
	"github.com/aquaproj/aqua/pkg/slsa"
	"github.com/aquaproj/aqua/pkg/unarchive"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/go-osenv/osenv"
	"net/http"
)

// Injectors from wire.go:

func InitializeListCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *list.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param, rt)
	slsaVerifier := slsa.New(downloader, fs)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	linker := link.New()
	checksumDownloader := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	unarchiver := unarchive.New()
	checker := policy.NewChecker()
	installpackageInstaller := installpackage.New(param, downloader, rt, fs, linker, executor, checksumDownloader, calculator, unarchiver, checker, verifier, slsaVerifier)
	controller := list.NewController(configFinder, configReader, installer, installpackageInstaller)
	return controller
}

func InitializeGenerateRegistryCommandController(ctx context.Context, param *config.Param, httpClient *http.Client) *genrgst.Controller {
	fs := afero.NewOsFs()
	repositoriesService := github.New(ctx)
	controller := genrgst.NewController(fs, repositoriesService)
	return controller
}

func InitializeInitCommandController(ctx context.Context, param *config.Param) *initcmd.Controller {
	repositoriesService := github.New(ctx)
	fs := afero.NewOsFs()
	controller := initcmd.New(repositoriesService, fs)
	return controller
}

func InitializeInitPolicyCommandController(ctx context.Context) *initpolicy.Controller {
	fs := afero.NewOsFs()
	controller := initpolicy.New(fs)
	return controller
}

func InitializeGenerateCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *generate.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param, rt)
	slsaVerifier := slsa.New(downloader, fs)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	fuzzyFinder := generate.NewFuzzyFinder()
	versionSelector := generate.NewVersionSelector()
	linker := link.New()
	checksumDownloader := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	unarchiver := unarchive.New()
	checker := policy.NewChecker()
	installpackageInstaller := installpackage.New(param, downloader, rt, fs, linker, executor, checksumDownloader, calculator, unarchiver, checker, verifier, slsaVerifier)
	controller := generate.New(configFinder, configReader, installer, repositoriesService, fs, fuzzyFinder, versionSelector, installpackageInstaller)
	return controller
}

func InitializeInstallCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *install.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param, rt)
	slsaVerifier := slsa.New(downloader, fs)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	linker := link.New()
	checksumDownloader := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	unarchiver := unarchive.New()
	checker := policy.NewChecker()
	installpackageInstaller := installpackage.New(param, downloader, rt, fs, linker, executor, checksumDownloader, calculator, unarchiver, checker, verifier, slsaVerifier)
	policyConfigReader := policy.NewConfigReader(fs)
	controller := install.New(param, configFinder, configReader, installer, installpackageInstaller, fs, rt, policyConfigReader, installpackageInstaller)
	return controller
}

func InitializeWhichCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *which.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param, rt)
	slsaVerifier := slsa.New(downloader, fs)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	osEnv := osenv.New()
	linker := link.New()
	checksumDownloader := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	unarchiver := unarchive.New()
	checker := policy.NewChecker()
	installpackageInstaller := installpackage.New(param, downloader, rt, fs, linker, executor, checksumDownloader, calculator, unarchiver, checker, verifier, slsaVerifier)
	controller := which.New(param, configFinder, configReader, installer, rt, osEnv, fs, linker, installpackageInstaller)
	return controller
}

func InitializeExecCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *exec2.Controller {
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	fs := afero.NewOsFs()
	linker := link.New()
	executor := exec.New()
	checksumDownloader := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	unarchiver := unarchive.New()
	checker := policy.NewChecker()
	verifier := cosign.NewVerifier(executor, fs, downloader, param, rt)
	slsaVerifier := slsa.New(downloader, fs)
	installer := installpackage.New(param, downloader, rt, fs, linker, executor, checksumDownloader, calculator, unarchiver, checker, verifier, slsaVerifier)
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	registryInstaller := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	osEnv := osenv.New()
	controller := which.New(param, configFinder, configReader, registryInstaller, rt, osEnv, fs, linker, installer)
	policyConfigReader := policy.NewConfigReader(fs)
	execController := exec2.New(installer, controller, executor, osEnv, fs, policyConfigReader, checker)
	return execController
}

func InitializeUpdateAquaCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *updateaqua.Controller {
	fs := afero.NewOsFs()
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	linker := link.New()
	executor := exec.New()
	checksumDownloader := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	unarchiver := unarchive.New()
	checker := policy.NewChecker()
	verifier := cosign.NewVerifier(executor, fs, downloader, param, rt)
	slsaVerifier := slsa.New(downloader, fs)
	installer := installpackage.New(param, downloader, rt, fs, linker, executor, checksumDownloader, calculator, unarchiver, checker, verifier, slsaVerifier)
	controller := updateaqua.New(param, fs, rt, repositoriesService, installer)
	return controller
}

func InitializeCopyCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *cp.Controller {
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	fs := afero.NewOsFs()
	linker := link.New()
	executor := exec.New()
	checksumDownloader := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	unarchiver := unarchive.New()
	checker := policy.NewChecker()
	verifier := cosign.NewVerifier(executor, fs, downloader, param, rt)
	slsaVerifier := slsa.New(downloader, fs)
	installer := installpackage.New(param, downloader, rt, fs, linker, executor, checksumDownloader, calculator, unarchiver, checker, verifier, slsaVerifier)
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	registryInstaller := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	osEnv := osenv.New()
	controller := which.New(param, configFinder, configReader, registryInstaller, rt, osEnv, fs, linker, installer)
	policyConfigReader := policy.NewConfigReader(fs)
	installController := install.New(param, configFinder, configReader, registryInstaller, installer, fs, rt, policyConfigReader, installer)
	cpController := cp.New(param, installer, fs, rt, controller, installController, policyConfigReader)
	return cpController
}

func InitializeUpdateChecksumCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *updatechecksum.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param, rt)
	slsaVerifier := slsa.New(downloader, fs)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	checksumDownloader := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	linker := link.New()
	calculator := checksum.NewCalculator()
	unarchiver := unarchive.New()
	checker := policy.NewChecker()
	installpackageInstaller := installpackage.New(param, downloader, rt, fs, linker, executor, checksumDownloader, calculator, unarchiver, checker, verifier, slsaVerifier)
	controller := updatechecksum.New(param, configFinder, configReader, installer, fs, rt, checksumDownloader, downloader, installpackageInstaller)
	return controller
}
