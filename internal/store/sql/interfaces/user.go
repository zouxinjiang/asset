package interfaces

import (
	sentity "github.com/zouxinjiang/axes/internal/store/sql/entity"
)

type (
	User interface {
		idResolver
		Create(username, displayName string, password []byte, mobile, email string, sourceID int64, originID string) (userID int64, err error)
		Update(userID int64, password []byte, displayName, mobile, email *string) (updated sentity.User, err error)
		DeleteByID(userID int64) (err error)
		List(filter sentity.Filter) (list []sentity.User, cnt int64, err error)
		ListResource(filter sentity.Filter) (list []sentity.User, cnt int64, err error)
		GetByID(userID int64) (sentity.User, error)
	}
	ResourceUserRel interface {
		resourcePropertyRel
		GetProperty(resourceID int64) (list []sentity.User, err error)
	}
)
