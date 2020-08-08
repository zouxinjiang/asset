package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	isql "github.com/zouxinjiang/axes/internal/store/sql/interfaces"
	"net/url"
	"time"
)

type (
	SqlTxStore struct {
		db *gorm.DB
	}

	Config struct {
		IsDebug           bool
		Host              string
		Port              uint16
		Username          string
		Password          string
		Schema            string
		DatabaseName      string
		ConnectionTimeout time.Duration
		ClientEncoding    string
		SslMode           string
	}
	SslMode string
)

const (
	SslModeDisable SslMode = "disable"
	SslModeAllow   SslMode = "allow"
	SslModePrefer  SslMode = "prefer"
	SslModeRequire SslMode = "require"
)

const (
	DriverPostgres isql.SqlStoreDriver = "postgres"
)

func (c Config) Driver() isql.SqlStoreDriver {
	return DriverPostgres
}

func (c Config) DataSourceName() string {
	//postgresql://[user[:password]@][netloc][:port][,...][/dbname][?param1=value1&...]
	dsn := fmt.Sprintf(`postgresql://%s:%s@%s:%d/%s`, c.Username, c.Password, c.Host, c.Port, c.DatabaseName)
	query := url.Values{}
	addFn := func(name string, value string) {
		if value != "" {
			query.Add(name, value)
		}
	}
	addFn("client_encoding", c.ClientEncoding)
	addFn("connect_timeout", fmt.Sprintf("%d", c.ConnectionTimeout/time.Millisecond))
	addFn("sslmode", c.SslMode)
	addFn("search_path", c.Schema)
	queryStr := query.Encode()
	if queryStr != "" {
		dsn += `?` + queryStr
	}
	return dsn
}

func (c Config) Connect() (*gorm.DB, error) {
	dsn := c.DataSourceName()
	conn, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	conn = conn.LogMode(c.IsDebug)
	return conn, nil
}

func (c Config) NewTxStore(db *gorm.DB) isql.SqlTxStore {
	s := SqlTxStore{
		db: db,
	}
	return s
}

func (s SqlTxStore) Resource() isql.Resource {
	return newResource(s.db)
}

func (s SqlTxStore) ResourceResourceRel() isql.ResourceResourceRel {
	return newResourceResourceRel(s.db)
}
