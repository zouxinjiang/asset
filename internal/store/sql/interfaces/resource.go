package interfaces

import (
	sentity "github.com/zouxinjiang/axes/internal/store/sql/entity"
)

type (
	idResolver interface {
		GetResourceIDMap(propertyIDs []int64) (propertyIDResourceIDMap map[int64]int64, err error)
		GetPropertyIDMap(resourceIDs []int64) (resourceIDPropertyIDMap map[int64]int64, err error)
	}
	resourcePropertyRel interface {
		BuildRelationship(resourceID, propertyID int64) error
		GetResource(propertyID int64) (list []sentity.Resource, err error)
		GetPropertyIDs(resourceID int64) (list []int64, err error)
	}

	// 资源表
	Resource interface {
		Create(typ sentity.ResourceType) (resourceID int64, err error)
		DeleteByID(resourceID int64) error
		GetByID(resourceID int64) (item sentity.Resource, err error)
		GetByType(typ sentity.ResourceType) (list []sentity.Resource, err error)
		UpdateByID(resourceID int64, typ *sentity.ResourceType) (updated sentity.Resource, err error)
	}

	// 资源-资源关系表
	ResourceResourceRel interface {
		GetDirectParents(nodeID int64, typ sentity.ResourceRelType) (parents []sentity.Resource, err error)
		GetDirectChildren(nodeID int64, typ sentity.ResourceRelType) (children []sentity.Resource, err error)
		GetAllParents(nodeID int64, typ sentity.ResourceRelType) (parents []sentity.Resource, err error)
		GetAllChildren(nodeID int64, typ sentity.ResourceRelType) (children []sentity.Resource, err error)

		Create(parentID, childID int64, typ sentity.ResourceRelType) (relID int64, err error)
		DeleteByID(relID int64) error
		DeleteByParentID(parentID int64, typ sentity.ResourceRelType) (rowAffect int64, err error)
		DeleteByChildID(childID int64, typ sentity.ResourceRelType) (rowAffect int64, err error)
		DeleteByParentIDAndChildID(parentID, childID int64, typ sentity.ResourceRelType) (rowAffect int64, err error)

		GetByID(ID int64) (e sentity.ResourceResourceRel, err error)
		GetByType(typ sentity.ResourceRelType) (list []sentity.ResourceResourceRel, err error)
		GetByParentIDAndChildID(parentID, childID int64, typ sentity.ResourceRelType) (e sentity.ResourceResourceRel, err error)
		GetByParentID(parentID int64, typ sentity.ResourceRelType) (list []sentity.ResourceResourceRel, err error)
		GetByChildID(childID int64, typ sentity.ResourceRelType) (list []sentity.ResourceResourceRel, err error)
	}
)
