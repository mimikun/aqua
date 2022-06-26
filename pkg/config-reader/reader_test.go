package reader_test

import (
	"testing"

	reader "github.com/clivm/clivm/pkg/config-reader"
	"github.com/clivm/clivm/pkg/config/aqua"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

func Test_configReader_Read(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name           string
		exp            *clivm.Config
		isErr          bool
		files          map[string]string
		configFilePath string
	}{
		{
			name:  "file isn't found",
			isErr: true,
		},
		{
			name: "normal",
			files: map[string]string{
				"clivm.yaml": `registries:
- type: standard
  ref: v2.5.0
packages:`,
			},
			configFilePath: "clivm.yaml",
			exp: &clivm.Config{
				Registries: clivm.Registries{
					"standard": {
						Type:      "github_content",
						Name:      "standard",
						Ref:       "v2.5.0",
						RepoOwner: "clivm",
						RepoName:  "clivm-registry",
						Path:      "registry.yaml",
					},
				},
				Packages: []*clivm.Package{},
			},
		},
		{
			name: "import package",
			files: map[string]string{
				"clivm.yaml": `registries:
- type: standard
  ref: v2.5.0
packages:
- name: suzuki-shunsuke/ci-info@v1.0.0
- import: clivm-installer.yaml
`,
				"clivm-installer.yaml": `packages:
- name: clivm/clivm-installer@v1.0.0
`,
			},
			configFilePath: "clivm.yaml",
			exp: &clivm.Config{
				Registries: clivm.Registries{
					"standard": {
						Type:      "github_content",
						Name:      "standard",
						Ref:       "v2.5.0",
						RepoOwner: "clivm",
						RepoName:  "clivm-registry",
						Path:      "registry.yaml",
					},
				},
				Packages: []*clivm.Package{
					{
						Name:     "suzuki-shunsuke/ci-info",
						Registry: "standard",
						Version:  "v1.0.0",
					},
					{
						Name:     "clivm/clivm-installer",
						Registry: "standard",
						Version:  "v1.0.0",
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()
			for name, body := range d.files {
				if err := afero.WriteFile(fs, name, []byte(body), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			reader := reader.New(fs)
			cfg := &clivm.Config{}
			if err := reader.Read(d.configFilePath, cfg); err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if diff := cmp.Diff(d.exp, cfg); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
