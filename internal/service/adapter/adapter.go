package adapter

import (
	"github.com/zouxinjiang/axes/internal/store"
)

type (
	Adapter struct {
		resource
	}
)

func New(store store.Store) *Adapter {
	return &Adapter{
		resource: resource{
			store: store,
		},
	}
}
