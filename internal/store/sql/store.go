package sql

import (
	"github.com/jinzhu/gorm"
	isql "github.com/zouxinjiang/axes/internal/store/sql/interfaces"
)

type (
	TransactionFunc func(txStore isql.SqlTxStore) error
	SqlStore        struct {
		db *gorm.DB
		isql.SqlTxStore
		cfg SqlStoreConfig
	}
	SqlStoreConfig interface {
		Driver() isql.SqlStoreDriver
		DataSourceName() string
		Connect() (*gorm.DB, error)
		NewTxStore(db *gorm.DB) isql.SqlTxStore
	}
)

func (s SqlStore) Transaction(fn TransactionFunc) error {
	db := s.db.Begin()
	tx := s.cfg.NewTxStore(db)
	err := fn(tx)
	if err != nil {
		db.Rollback()
		return err
	} else {
		db.Commit()
	}
	return nil
}

func NewSqlStore(config SqlStoreConfig) (*SqlStore, error) {
	db, err := config.Connect()
	if err != nil {
		return nil, err
	}
	s := &SqlStore{
		db: db,
	}
	return s, nil
}
