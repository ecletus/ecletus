package aghape

import (
	"sync"

	"github.com/aghape/aghape/gid"
)

var instances = &aghapeInstances{}

type aghapeInstances struct {
	data sync.Map
}

func (ai *aghapeInstances) with(agp *Aghape) func() {
	ai.data.Store(gid.GID(), agp)
	return func() {
		ai.data.Delete(gid.GID())
	}
}

func Get() *Aghape {
	agp, _ := instances.data.Load(gid.GID())
	return agp.(*Aghape)
}
