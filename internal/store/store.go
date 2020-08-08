package store

import (
	"github.com/zouxinjiang/axes/internal/store/sql"
)

type (
	Store struct {
		SqlStore *sql.SqlStore
	}
	Option func(s *Store) error
)

func NewStore(options ...Option) (*Store, error) {
	s := &Store{}
	for _, v := range options {
		err := v(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func WithSqlStore(cfg sql.SqlStoreConfig) Option {
	return func(s *Store) error {
		sqlStore, err := sql.NewSqlStore(cfg)
		if err != nil {
			return err
		}
		s.SqlStore = sqlStore
		return nil
	}
}
