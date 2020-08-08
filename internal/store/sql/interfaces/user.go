package interfaces

import (
	sentity "github.com/zouxinjiang/axes/internal/store/sql/entity"
)

type (
	User interface {
		idResolver
		Create() (sourceID int64, err error)
		Update() (updated sentity.User, err error)
		DeleteByID(sourceID int64) (err error)
		List(filter sentity.Filter) (list []sentity.User, cnt int64, err error)
		ListResource(filter sentity.Filter) (list []sentity.User, cnt int64, err error)
		GetByID(sourceID int64) (sentity.User, error)
	}
	ResourceUserRel interface {
		resourcePropertyRel
		GetProperty(resourceID int64) (list []sentity.User, err error)
	}
)
