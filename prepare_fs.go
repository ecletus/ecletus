// +build !bindata

package ecletus

import (
	"github.com/moisespsena/go-assetfs"
	"github.com/moisespsena/go-assetfs/assetfsapi"
)

func (a *Ecletus) preparePlugins() {
	plugins := a.Plugins()
	if _, ok := plugins.FS().(*assetfs.AssetFileSystem); ok {
		plugins.SetAssetFSPathRegister(func(fs assetfsapi.PathRegistrator, pth string) error {
			return fs.PrependPath(pth)
		})
	}
}
