package service

import (
	"github.com/zouxinjiang/axes/internal/store"
)

type (
	Service struct {
		store *store.Store
	}
)

func NewService(store *store.Store) *Service {
	s := &Service{
		store: store,
	}

	return s
}
