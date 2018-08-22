package plugin

import (
	"plugin"

	"github.com/aghape/aghape"
	"github.com/go-errors/errors"
	"github.com/moisespsena/go-error-wrap"
)

var ErrInvalidApplication = errors.New("Invalid Application")

func LoadGoPlugin(pth string) (*aghape.Aghape, error) {
	plug, err := plugin.Open(pth)
	if err != nil {
		return nil, err
	}

	symGreeter, err := plug.Lookup("Application")
	if err != nil {
		return nil, errwrap.Wrap(err, "Application Lookup")
	}

	var app aghape.ApplicationInterface
	app, ok := symGreeter.(aghape.ApplicationInterface)
	if !ok {
		return nil, ErrInvalidApplication
	}

	return aghape.LoadApplication(app)
}
