package ecletus

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/moisespsena-go/pluggable"

	logging_helpers "github.com/moisespsena-go/logging-helpers"

	"github.com/moisespsena-go/logging"

	"github.com/go-errors/errors"

	"github.com/moisespsena-go/task"

	defaultlogger "github.com/moisespsena-go/default-logger"

	"github.com/ecletus/container"
	"github.com/ecletus/plug"
	"github.com/ecletus/sites"
	"github.com/moisespsena-go/assetfs"
	"github.com/moisespsena-go/assetfs/assetfsapi"
	errwrap "github.com/moisespsena-go/error-wrap"
	path_helpers "github.com/moisespsena-go/path-helpers"

	"github.com/ecletus/core"
)

const (
	ECLETUS      = "ecletus"
	SITES_CONFIG = "ecletus:SitesConfig"
	SETUP_CONFIG = "ecletus:SetupConfig"
	CONTAINER    = "ecletus:Container"
	ASSETFS      = "ecletus:Assetfs"
	CONFIG_DIR   = "ecletus:ConfigDir"

	DEFAULT_CONFIG_DIR = "config"
)

var log = defaultlogger.GetOrCreateLogger(path_helpers.GetCalledDir())

type Ecletus struct {
	task.Tasks
	AppName                     string
	ConfigDir                   *ConfigDir
	SitesConfig                 *sites.Config
	SetupConfig                 *core.SetupConfig
	AssetFS                     assetfsapi.Interface
	PubicFS                     *assetfs.AssetFileSystem
	TempFS                      *assetfs.AssetFileSystem
	Container                   *container.Container
	plugins                     []interface{}
	PrePluginsRegisterCallbacks []func(ecl *Ecletus) error
	PreInitCallbacks            []func(ecl *Ecletus) error
	done                        []func()
	Stderr                      io.Writer
	BasicSystemInfo             *BasicSystemInfo
}

func (this *Ecletus) PrePluginsRegister(f ...func(ecl *Ecletus) error) *Ecletus {
	this.PrePluginsRegisterCallbacks = append(this.PrePluginsRegisterCallbacks, f...)
	return this
}

func (this *Ecletus) PreInit(f ...func(ecl *Ecletus) error) *Ecletus {
	this.PreInitCallbacks = append(this.PreInitCallbacks, f...)
	return this
}

func (this *Ecletus) Plugins() *plug.Plugins {
	return this.Container.Plugins
}

func (this *Ecletus) Options() *plug.Options {
	return this.Container.Options
}

func (this *Ecletus) Done(f ...func()) {
	this.done = append(this.done, f...)
}

func (this *Ecletus) LoadLogLevels() {
	var cfg logging_helpers.LoggingConfig
	err := this.ConfigDir.Load(&cfg, "log.yaml", "log.yml")
	defaultLevel := cfg.GetLevel()

	if err == nil {
		for _, mod := range cfg.Modules {
			if mod.Name != "" {
				log := logging.GetOrCreateLogger(mod.Name)
				logging.SetLogLevel(log, mod.GetLevel(defaultLevel), mod.Name)
				if backends := mod.Backend(); len(backends) > 0 {
					func(backends ...logging.BackendCloser) {
						var bce = make([]logging.Backend, len(backends))
						for i, bc := range backends {
							bce[i] = bc
						}
						this.Done(func() {
							for _, bce := range backends {
								bce.Close()
							}
						})
						log.SetBackend(logging.MultiLogger(bce...))
					}(backends...)
				}
			}
		}
	} else {
		if !os.IsNotExist(err) {
			panic(errors.New("Ecletus.LoadLogLevels: " + err.Error()))
		}
	}
}

func (this *Ecletus) Init(plugins []interface{}) error {
	if this.AppName == "" {
		this.AppName = os.Args[0]
	}

	this.plugins = plugins
	if this.ConfigDir == nil {
		this.ConfigDir = NewConfigDir(this.AppName)
	}

	if this.SitesConfig == nil {
		this.SitesConfig = NewSitesConfig(this.ConfigDir)
	}

	if this.SetupConfig == nil {
		this.SetupConfig = NewSetupConfig(this.SitesConfig)
	}

	if this.Container == nil {
		pls := plug.New(this.AssetFS)
		this.Container = container.New(pls)
	}

	this.preparePlugins()

	options := this.Container.Options
	options.Set(ECLETUS, this)
	options.Set(SITES_CONFIG, this.SitesConfig)
	options.Set(SETUP_CONFIG, this.SetupConfig)
	options.Set(CONTAINER, this.Container)
	options.Set(ASSETFS, this.AssetFS)
	options.Set(CONFIG_DIR, this.ConfigDir)

	this.LoadLogLevels()

	for _, f := range this.PrePluginsRegisterCallbacks {
		if err := f(this); err != nil {
			return errwrap.Wrap(err, "Pre Plugins")
		}
	}

	if err := this.Container.Plugins.Add(this.plugins...); err != nil {
		return errwrap.Wrap(err, "Register plugins")
	}

	// temp fs
	tmpDir := this.SetupConfig.TempDir()
	this.TempFS = assetfs.NewAssetFileSystem()
	_ = this.TempFS.RegisterPath(tmpDir, false)

	// public fs
	publicDir := path.Join(this.SetupConfig.Root(), "public")
	this.PubicFS = assetfs.NewAssetFileSystem()
	_ = this.PubicFS.RegisterPath(publicDir, false)

	for _, f := range this.PreInitCallbacks {
		if err := f(this); err != nil {
			return errwrap.Wrap(err, "Pre Init")
		}
	}

	return this.Container.Init()
}

func (this *Ecletus) Setup(ta task.Appender) (err error) {
	defer instances.with(this)()
	return this.Tasks.Setup(ta)
}

func (this *Ecletus) Run() (err error) {
	defer instances.with(this)()
	defer func() {
		for _, done := range this.done {
			done()
		}
	}()
	return this.Tasks.Run()
}

func (this *Ecletus) Start(done func()) (stop task.Stoper, err error) {
	ldone := instances.with(this)
	return this.Tasks.Start(func() {
		defer func() {
			ldone()
			done()
			for _, done := range this.done {
				done()
			}
		}()
		log.Info("done.")
	})
}

func (this *Ecletus) Migrate(ctx context.Context) error {
	return this.Container.Migrate(ctx)
}

func (this *Ecletus) Main(main func()) {
	main()
}

func New() *Ecletus {
	a := &Ecletus{}
	a.Tasks.SetLog(log)
	return a
}

func FromOptions(options pluggable.Options) *Ecletus {
	return options.GetInterface(ECLETUS).(*Ecletus)
}
