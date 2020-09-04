package ecletus

import (
	"os"

	"github.com/moisespsena-go/maps"

	"github.com/ecletus/core"
	"github.com/ecletus/sites"
	"github.com/ecletus/sites/dir_config"
)

var (
	Root, _ = os.Getwd()
	HOME    = os.Getenv("HOME")
)

func NewSitesConfig(configDir *ConfigDir) *sites.Config {
	root := configDir.Path("sites")
	cfg, err := dir_config.LoadMainConfig(root, func(dir, name string, isdir bool) string {
		if isdir && dir == root {
			switch name {
			case "site", "_template":
				return "site_template"
			}
		}
		return ""
	})
	if err != nil {
		panic(err)
	}

	Config := &sites.Config{}
	if err = cfg.CopyTo(Config); err != nil {
		panic(err)
	}
	if Config.DataDir == "" {
		Config.DataDir = "data"
	}
	Config.Raw = cfg
	if rawSite, ok := cfg["site_template"]; ok {
		Config.SiteTemplate.Raw = rawSite.(maps.MapSI)
	}
	return Config
}

func NewSetupConfig(sitesConfig *sites.Config) *core.SetupConfig {
	return core.Setup(core.SetupOptions{
		Home: HOME,
	})
}
