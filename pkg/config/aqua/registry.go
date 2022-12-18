package aqua

import (
	"fmt"
	"path/filepath"

	"github.com/aquaproj/aqua/pkg/config/registry"
	"github.com/aquaproj/aqua/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Registry struct {
	Name           string                   `validate:"required" json:"name,omitempty"`
	Type           string                   `validate:"required" json:"type,omitempty" jsonschema:"enum=standard,enum=local,enum=github_content"`
	RepoOwner      string                   `yaml:"repo_owner" json:"repo_owner,omitempty"`
	RepoName       string                   `yaml:"repo_name" json:"repo_name,omitempty"`
	Ref            string                   `json:"ref,omitempty"`
	Path           string                   `validate:"required" json:"path,omitempty"`
	Cosign         *registry.Cosign         `json:"cosign,omitempty"`
	SLSAProvenance *registry.SLSAProvenance `json:"slsa_provenance,omitempty" yaml:"slsa_provenance"`
}

const (
	RegistryTypeGitHubContent = "github_content"
	RegistryTypeLocal         = "local"
	RegistryTypeStandard      = "standard"
)

func (registry *Registry) Validate() error {
	switch registry.Type {
	case RegistryTypeLocal:
		return registry.validateLocal()
	case RegistryTypeGitHubContent:
		return registry.validateGitHubContent()
	default:
		return logerr.WithFields(errInvalidRegistryType, logrus.Fields{ //nolint:wrapcheck
			"registry_type": registry.Type,
		})
	}
}

func (registry *Registry) validateLocal() error {
	if registry.Path == "" {
		return errPathIsRequired
	}
	return nil
}

func (registry *Registry) validateGitHubContent() error {
	if registry.RepoOwner == "" {
		return errRepoOwnerIsRequired
	}
	if registry.RepoName == "" {
		return errRepoNameIsRequired
	}
	if registry.Ref == "" {
		return errRefIsRequired
	}
	return nil
}

func (registry *Registry) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias Registry
	a := alias(*registry)
	if err := unmarshal(&a); err != nil {
		return err
	}
	if a.Type == RegistryTypeStandard {
		a.Type = RegistryTypeGitHubContent
		if a.Name == "" {
			a.Name = RegistryTypeStandard
		}
		if a.RepoOwner == "" {
			a.RepoOwner = "aquaproj"
		}
		if a.RepoName == "" {
			a.RepoName = "aqua-registry"
		}
		if a.Path == "" {
			a.Path = "registry.yaml"
		}
	}
	*registry = Registry(a)
	return nil
}

func (registry *Registry) GetFilePath(rootDir, cfgFilePath string) (string, error) {
	switch registry.Type {
	case RegistryTypeLocal:
		return util.Abs(filepath.Dir(cfgFilePath), registry.Path), nil
	case RegistryTypeGitHubContent:
		return filepath.Join(rootDir, "registries", registry.Type, "github.com", registry.RepoOwner, registry.RepoName, registry.Ref, registry.Path), nil
	}
	return "", errInvalidRegistryType
}

func (registry *Registry) SLSASourceURI() string {
	sp := registry.SLSAProvenance
	if sp == nil {
		return ""
	}
	if sp.SourceURI != nil {
		return *sp.SourceURI
	}
	repoOwner := sp.RepoOwner
	repoName := sp.RepoName
	if repoOwner == "" {
		repoOwner = registry.RepoOwner
	}
	if repoName == "" {
		repoName = registry.RepoName
	}
	return fmt.Sprintf("github.com/%s/%s", repoOwner, repoName)
}
