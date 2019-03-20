package plugin

import (
	"plugin"

	"github.com/ecletus/ecletus"
	"github.com/go-errors/errors"
	"github.com/moisespsena/go-error-wrap"
)

var ErrInvalidApplication = errors.New("Invalid Application")

func LoadGoPlugin(pth string) (*ecletus.Ecletus, error) {
	plug, err := plugin.Open(pth)
	if err != nil {
		return nil, err
	}

	symGreeter, err := plug.Lookup("Application")
	if err != nil {
		return nil, errwrap.Wrap(err, "Application Lookup")
	}

	var app ecletus.ApplicationInterface
	app, ok := symGreeter.(ecletus.ApplicationInterface)
	if !ok {
		return nil, ErrInvalidApplication
	}

	return ecletus.LoadApplication(app)
}
