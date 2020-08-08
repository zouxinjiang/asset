package entity

import (
	sentity "github.com/zouxinjiang/axes/internal/store/sql/entity"
)

type (
	Filter = sentity.Filter
)

const (
	Eq        = sentity.Eq
	Ne        = sentity.Ne
	Gt        = sentity.Gt
	Ge        = sentity.Ge
	Lt        = sentity.Lt
	Le        = sentity.Le
	Like      = sentity.Like
	NotLike   = sentity.NotLike
	IsNull    = sentity.IsNull
	IsNotNull = sentity.IsNull
	In        = sentity.In
	NotIn     = sentity.NotIn
)
