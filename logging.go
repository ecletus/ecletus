package aghape

import (
	"strings"

	"github.com/op/go-logging"
)

type LogLevel struct {
	Level string
}

var levels = map[string]logging.Level{
	"CRITICAL": logging.CRITICAL,
	"C":        logging.CRITICAL,
	"ERROR":    logging.ERROR,
	"E":        logging.ERROR,
	"WARNING":  logging.WARNING,
	"W":        logging.WARNING,
	"NOTICE":   logging.NOTICE,
	"N":        logging.NOTICE,
	"INFO":     logging.INFO,
	"I":        logging.INFO,
	"DEBUG":    logging.DEBUG,
	"D":        logging.DEBUG,
}

func (ll LogLevel) GetLevel(defaul ...logging.Level) logging.Level {
	if l, ok := levels[strings.ToUpper(ll.Level)]; ok {
		return l
	}
	for _, d := range defaul {
		return d
	}
	return logging.DEBUG
}

type ModuleLoggingConfig struct {
	LogLevel `yaml:",inline"`
	Name    string
	Dst     string
	Options map[string]interface{}
}

type LoggingConfig struct {
	LogLevel `yaml:",inline"`
	Modules []ModuleLoggingConfig
}
