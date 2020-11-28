package entity

import (
	"time"
)

type (
	User struct {
		ID           int64
		UserName     string
		DisplayName  *string
		Password     []byte
		Mobile       *string
		Email        *string
		UserSourceID int64
		OriginID     *string
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}
)
