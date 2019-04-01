package ecletus

import (
	"github.com/moisespsena/go-assetfs/assetfsapi"
	"github.com/moisespsena-go/error-wrap"
)

type ApplicationInterface interface {
	GetAssetFS() assetfsapi.Interface
	GetPlugins() []interface{}
	PreInit(a *Ecletus) error
	PrePlugins(a *Ecletus) error
	PostInit(a *Ecletus) error
}

type Application struct {
	AssetFS         assetfsapi.Interface
	Plugins         []interface{}
	PrePluginstFunc func(a *Ecletus) error
	PreInitFunc     func(a *Ecletus) error
	PostInitFunc    func(a *Ecletus) error
}

func (app *Application) GetAssetFS() assetfsapi.Interface {
	return app.AssetFS
}

func (app *Application) GetPlugins() []interface{} {
	return app.Plugins
}

func (app *Application) PreInit(a *Ecletus) error {
	if app.PreInitFunc == nil {
		return nil
	}
	return app.PreInitFunc(a)
}

func (app *Application) PrePlugins(a *Ecletus) error {
	if app.PrePluginstFunc == nil {
		return nil
	}
	return app.PrePluginstFunc(a)
}

func (app *Application) PostInit(a *Ecletus) error {
	if app.PostInitFunc == nil {
		return nil
	}
	return app.PostInitFunc(a)
}

func LoadApplication(app ApplicationInterface) (*Ecletus, error) {
	agp := &Ecletus{
		AssetFS:    app.GetAssetFS(),
		PreInit:    app.PreInit,
		PrePlugins: app.PrePlugins,
	}

	if err := agp.Init(app.GetPlugins()); err != nil {
		return nil, errwrap.Wrap(err, "Init")
	}

	if err := app.PostInit(agp); err != nil {
		return nil, errwrap.Wrap(err, "Post Init")
	}

	return agp, nil
}
