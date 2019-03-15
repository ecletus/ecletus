package aghape

import (
	"os"
	"path/filepath"

	"strings"

	"regexp"

	"github.com/jinzhu/configor"
)

type ConfigDir struct {
	dir      string
	Configor *configor.Configor
}

func NewConfigDir(envPrefix ...string) *ConfigDir {
	prefix := strings.Trim(
		regexp.MustCompile(`_{2,}`).ReplaceAllString(
			regexp.MustCompile(`[-\.]`).ReplaceAllString(envPrefix[0], "_"),
			"_"),
		"_")

	dir := os.Getenv(prefix + "_CONFIG_DIR")
	if dir == "" {
		dir = DEFAULT_CONFIG_DIR
	}
	if len(envPrefix) == 0 {
		envPrefix = []string{os.Args[0]}
	}

	return &ConfigDir{dir, configor.New(&configor.Config{
		Debug:     false,
		ENVPrefix: prefix,
		Verbose:   false,
	})}
}

func (c *ConfigDir) Load(config interface{}, files ...string) error {
	dir := c.dir
	var news []string
	for i, f := range files {
		files[i] = filepath.Join(dir, f)
		if strings.HasSuffix(f, ".yml") {
			news = append(news, filepath.Join(dir, strings.TrimSuffix(f, "yml")+"yaml"))
		} else if strings.HasSuffix(f, ".yaml") {
			news = append(news, filepath.Join(dir, strings.TrimSuffix(f, "yaml")+"yml"))
		}
	}
	return configor.Load(config, append(files, news...)...)
}

func (c *ConfigDir) Path(pth ...string) string {
	return filepath.Join(append([]string{c.dir}, pth...)...)
}

func (c *ConfigDir) Exists(pth ...string) (ok bool, err error) {
	if _, err = os.Stat(c.Path(pth...)); err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}
		return
	}
	return true, nil
}

func (c *ConfigDir) Paths(pth ...string) []string {
	dir := c.dir
	for i, f := range pth {
		pth[i] = filepath.Join(dir, f)
	}
	return pth
}

func (c *ConfigDir) Dir() string {
	return c.dir
}
