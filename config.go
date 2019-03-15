package aghape

import (
	"os"
	"strings"

	"github.com/aghape/core"
	"github.com/aghape/sites"
)

var (
	Root, _ = os.Getwd()
	HOME    = os.Getenv("HOME")
)

func NewSitesConfig(configDir *ConfigDir) *sites.Config {
	Config := &sites.Config{}

	if err := configDir.Load(Config, "database.yml", "smtp.yml", "application.yml", "sites.yml"); err != nil {
		panic(err)
	}

	if Config.Prefix != "" {
		Config.Prefix = "/" + strings.Trim(Config.Prefix, "/") + "/"
	}
	return Config
}

func NewSetupConfig(sitesConfig *sites.Config) *core.SetupConfig {
	return core.Setup(core.SetupOptions{
		Home:   HOME,
		Prefix: sitesConfig.Prefix,
	})
}
