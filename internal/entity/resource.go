package entity

import (
	sentity "github.com/zouxinjiang/axes/internal/store/sql/entity"
)

type (
	ResourceType    = sentity.ResourceType
	ResourceRelType = sentity.ResourceRelType

	ResourceCreator interface {
		ResourceType() ResourceType
		SetID(id int64)
		GetID() int64
		GetPropertyID() int64
		SetPropertyID(id int64)
	}

	property struct {
		ID         int64
		PropertyID int64
	}
)

func (p property) GetPropertyID() int64 {
	return p.PropertyID
}

func (p *property) SetPropertyID(id int64) {
	p.PropertyID = id
}

func (p property) GetID() int64 {
	return p.ID
}

func (p *property) SetID(id int64) {
	p.ID = id
}

func GenRelType(typ1, typ2 ResourceType) ResourceRelType {
	return ResourceRelType(typ1.String() + "-" + typ2.String())
}
