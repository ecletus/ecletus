package ecletus

import (
	"io"
	"os"
	"path"

	"github.com/go-errors/errors"
	"github.com/op/go-logging"

	"github.com/moisespsena-go/task"

	"github.com/moisespsena-go/default-logger"

	"github.com/ecletus/cli"
	"github.com/ecletus/container"
	"github.com/ecletus/core"
	"github.com/ecletus/plug"
	"github.com/ecletus/sites"
	"github.com/moisespsena/go-assetfs"
	"github.com/moisespsena/go-assetfs/assetfsapi"
	"github.com/moisespsena/go-error-wrap"
	"github.com/moisespsena-go/path-helpers"
)

const (
	AGHAPE       = "ecletus"
	SITES_CONFIG = "ecletus:SitesConfig"
	SETUP_CONFIG = "ecletus:SetupConfig"
	CONTAINER    = "ecletus:Container"
	ASSETFS      = "ecletus:Assetfs"
	CONFIG_DIR   = "ecletus:ConfigDir"

	DEFAULT_CONFIG_DIR = "config"
)

var log = defaultlogger.NewLogger(path_helpers.GetCalledDir())

type Ecletus struct {
	task.Tasks
	AppName     string
	ConfigDir   *ConfigDir
	SitesConfig *sites.Config
	SetupConfig *core.SetupConfig
	AssetFS     assetfsapi.Interface
	PubicFS     *assetfs.AssetFileSystem
	TempFS      *assetfs.AssetFileSystem
	Container   *container.Container
	plugins     []interface{}
	PrePlugins  func(a *Ecletus) error
	PreInit     func(a *Ecletus) error
	done        []func()
	cli         *cli.CLI
	Stderr      io.Writer
}

func (a *Ecletus) Plugins() *plug.Plugins {
	return a.Container.Plugins
}

func (a *Ecletus) Options() *plug.Options {
	return a.Container.Options
}

func (a *Ecletus) Done(f ...func()) {
	a.done = append(a.done, f...)
}

func (a *Ecletus) LoadLogLevels() {
	var cfg LoggingConfig
	err := a.ConfigDir.Load(&cfg, "log.yaml", "log.yml")
	defaultLevel := cfg.GetLevel()
	if err == nil {
		for _, mod := range cfg.Modules {
			if mod.Name != "" {
				logging.SetLevel(mod.GetLevel(defaultLevel), mod.Name)
			}
		}
	} else {
		if !os.IsNotExist(err) {
			panic(errors.New("Ecletus.LoadLogLevels: " + err.Error()))
		}
	}
}

func (a *Ecletus) Init(plugins []interface{}) error {
	if a.AppName == "" {
		a.AppName = os.Args[0]
	}

	a.plugins = plugins
	if a.ConfigDir == nil {
		a.ConfigDir = NewConfigDir(a.AppName)
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
	options.Set(AGHAPE, a)
	options.Set(SITES_CONFIG, a.SitesConfig)
	options.Set(SETUP_CONFIG, a.SetupConfig)
	options.Set(CONTAINER, a.Container)
	options.Set(ASSETFS, a.AssetFS)
	options.Set(CONFIG_DIR, a.ConfigDir)

	a.LoadLogLevels()

	if a.PrePlugins != nil {
		if err := a.PrePlugins(a); err != nil {
			return errwrap.Wrap(err, "Pre Plugins")
		}
	}

	if err := a.Container.Plugins.Add(a.plugins...); err != nil {
		return errwrap.Wrap(err, "Register plugins")
	}

	// temp fs
	tmpDir := a.SetupConfig.TempDir()
	a.TempFS = assetfs.NewAssetFileSystem()
	_ = a.TempFS.RegisterPath(tmpDir, false)

	// public fs
	publicDir := path.Join(a.SetupConfig.Root(), "public")
	a.PubicFS = assetfs.NewAssetFileSystem()
	_ = a.PubicFS.RegisterPath(publicDir, false)

	if a.PreInit != nil {
		if err := a.PreInit(a); err != nil {
			return errwrap.Wrap(err, "Pre Init")
		}
	}

	return a.Container.Init()
}

func (a *Ecletus) Setup(ta task.Appender) (err error) {
	defer instances.with(a)()
	if err = a.CLI().Execute(); err != nil {
		return
	}

	return a.Tasks.Setup(ta)
}

func (a *Ecletus) Run() (err error) {
	defer instances.with(a)()
	defer func() {
		for _, done := range a.done {
			done()
		}
	}()
	return a.Tasks.Run()
}

func (a *Ecletus) Start(done func()) (stop task.Stoper, err error) {
	ldone := instances.with(a)
	return a.Tasks.Start(func() {
		defer func() {
			ldone()
			done()
			for _, done := range a.done {
				done()
			}
		}()
		log.Info("done.")
	})
}

func (a *Ecletus) Migrate() error {
	return a.Container.Migrate()
}

func (a *Ecletus) Main(main func()) {
	main()
}

func (a *Ecletus) CLI() *cli.CLI {
	if a.cli == nil {
		a.cli = a.Container.CLI()
		a.cli.Stderr = a.Stderr
	}
	return a.cli
}

func New() *Ecletus {
	a := &Ecletus{}
	a.Tasks.SetLog(log)
	return a
}
