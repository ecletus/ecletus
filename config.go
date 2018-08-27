package aghape

import (
	"os"
	"strings"

	"github.com/aghape/core"
	"github.com/aghape/sites"
	"github.com/jinzhu/configor"
)

var (
	Root, _   = os.Getwd()
	ConfigDir string
	HOME      = os.Getenv("HOME")
)

func NewSitesConfig(configDir string) *sites.Config {
	ConfigDir = os.Getenv("CONFIG_DIR")
	if ConfigDir == "" {
		ConfigDir = configDir
	}

	Config := &sites.Config{}

	if err := configor.Load(Config, ConfigDir+"/database.yml", ConfigDir+"/smtp.yml", ConfigDir+"/application.yml", ConfigDir+"/sites.yml"); err != nil {
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
