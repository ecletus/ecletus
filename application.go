package aghape

import (
	"github.com/moisespsena/go-assetfs/api"
	"github.com/moisespsena/go-error-wrap"
)

type ApplicationInterface interface {
	GetAssetFS() api.Interface
	GetPlugins() []interface{}
	PreInit(a *Aghape) error
	PrePlugins(a *Aghape) error
	PostInit(a *Aghape) error
}

type Application struct {
	AssetFS         api.Interface
	Plugins         []interface{}
	PrePluginstFunc func(a *Aghape) error
	PreInitFunc     func(a *Aghape) error
	PostInitFunc    func(a *Aghape) error
}

func (app *Application) GetAssetFS() api.Interface {
	return app.AssetFS
}

func (app *Application) GetPlugins() []interface{} {
	return app.Plugins
}

func (app *Application) PreInit(a *Aghape) error {
	if app.PreInitFunc == nil {
		return nil
	}
	return app.PreInitFunc(a)
}

func (app *Application) PrePlugins(a *Aghape) error {
	if app.PrePluginstFunc == nil {
		return nil
	}
	return app.PrePluginstFunc(a)
}

func (app *Application) PostInit(a *Aghape) error {
	if app.PostInitFunc == nil {
		return nil
	}
	return app.PostInitFunc(a)
}

func LoadApplication(app ApplicationInterface) (*Aghape, error) {
	agp := &Aghape{
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
