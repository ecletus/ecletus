package ecletus

import (
	"sync"

	"github.com/ecletus/ecletus/gid"
)

var instances = &ecletusInstances{}

type ecletusInstances struct {
	data sync.Map
}

func (ai *ecletusInstances) with(agp *Ecletus) func() {
	ai.data.Store(gid.GID(), agp)
	return func() {
		ai.data.Delete(gid.GID())
	}
}

func Get() *Ecletus {
	agp, _ := instances.data.Load(gid.GID())
	return agp.(*Ecletus)
}
