package aghape

import (
	"path"

	"github.com/aghape/container"
	"github.com/aghape/core"
	"github.com/aghape/plug"
	"github.com/aghape/sites"
	"github.com/moisespsena/go-assetfs"
	"github.com/moisespsena/go-assetfs/api"
	"github.com/moisespsena/go-error-wrap"
	"github.com/moisespsena/go-path-helpers"
)

const (
	SITES_CONFIG       = "aghape:SitesConfig"
	SETUP_CONFIG       = "aghape:SetupConfig"
	CONTAINER          = "aghape:Container"
	ASSETFS            = "aghape:Assetfs"
	DEFAULT_CONFIG_DIR = "config"
)

type Aghape struct {
	ConfigDir   string
	SitesConfig *sites.Config
	SetupConfig *core.SetupConfig
	AssetFS     api.Interface
	Container   *container.Container
	plugins     []interface{}
	PrePlugins  func(a *Aghape) error
	PreInit     func(a *Aghape) error
}

func (a *Aghape) Plugins() *plug.Plugins {
	return a.Container.Plugins
}

func (a *Aghape) Options() *plug.Options {
	return a.Container.Options
}

func (a *Aghape) Init(plugins []interface{}) error {
	a.plugins = plugins
	if a.ConfigDir == "" {
		a.ConfigDir = DEFAULT_CONFIG_DIR
	}

	if a.SitesConfig == nil {
		a.SitesConfig = NewSitesConfig(a.ConfigDir)
	}

	if a.SetupConfig == nil {
		a.SetupConfig = NewSetupConfig(a.SitesConfig)
	}

	if a.Container == nil {
		pls := plug.New(a.AssetFS)
		a.Container = container.New(pls)
	}

	a.preparePlugins()

	options := a.Container.Options
	options.Set(SITES_CONFIG, a.SitesConfig)
	options.Set(SETUP_CONFIG, a.SetupConfig)
	options.Set(CONTAINER, a.Container)
	options.Set(ASSETFS, a.AssetFS)

	if a.PrePlugins != nil {
		if err := a.PrePlugins(a); err != nil {
			return errwrap.Wrap(err, "Pre Plugins")
		}
	}

	a.Container.Plugins.Add(a.plugins...)

	publicDir := path.Join(a.SetupConfig.Root(), "public")

	if path_helpers.IsExistingDir(publicDir) {
		a.AssetFS.RegisterPath(publicDir)
	}

	tmpDir := a.SetupConfig.TempDir()
	tmpFS := assetfs.NewAssetFileSystem()
	tmpFS.RegisterPath(tmpDir)
	a.AssetFS.NameSpace("tmp").Provider(tmpFS)

	if a.PreInit != nil {
		if err := a.PreInit(a); err != nil {
			return errwrap.Wrap(err, "Pre Init")
		}
	}

	return a.Container.Init()
}

func (a *Aghape) Migrate() error {
	return a.Container.Migrate()
}

func (a *Aghape) Execute() error {
	return a.Container.CLI().Execute()
}

func (a *Aghape) ExecuteAlone() {
	a.Container.CLI().ExecuteAlone()
}
