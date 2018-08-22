// +build !bindata

package aghape

import (
	"github.com/moisespsena/go-assetfs"
	"github.com/moisespsena/go-assetfs/api"
)

func (a *Aghape) preparePlugins() {
	plugins := a.Plugins()
	if _, ok := plugins.AssetFS.(*assetfs.AssetFileSystem); ok {
		plugins.SetAssetFSPathRegister(func(fs api.Interface, pth string) error {
			return fs.PrependPath(pth)
		})
	}
}
