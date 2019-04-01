package ecletus

import (
	"os"

	"github.com/ecletus/core"
	"github.com/ecletus/sites"
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
	return Config
}

func NewSetupConfig(sitesConfig *sites.Config) *core.SetupConfig {
	return core.Setup(core.SetupOptions{
		Home: HOME,
	})
}
