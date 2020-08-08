package adapter

import (
	"fmt"
	"github.com/zouxinjiang/axes/internal/entity"
	"github.com/zouxinjiang/axes/internal/store"
	isql "github.com/zouxinjiang/axes/internal/store/sql/interfaces"
	"github.com/zouxinjiang/axes/pkg/errors"
	"github.com/zouxinjiang/axes/pkg/tools"
)

type (
	resource struct {
		store store.Store
	}

	PropertyCreator    func(tx isql.SqlTxStore, resourceID int64) (propertyID int64, err error)
	PropertyUpdater    func(tx isql.SqlTxStore, resourceID int64) (updated entity.ResourceCreator, err error)
	PropertyGetter     func(tx isql.SqlTxStore, resourceID int64) (info entity.ResourceCreator, err error)
	PropertyDeleter    func(tx isql.SqlTxStore, resourceID int64) error
	PropertyLister     func(tx isql.SqlTxStore, filter entity.Filter) (resources []entity.ResourceCreator, count int64, getResourceMap GetResourceIDsMap, err error)
	PropertyCheckExist func(tx isql.SqlTxStore, withPatentCheck bool, nodeID, parentID int64, parentTyp entity.ResourceType) error
	GetResourceIDsMap  func(propertyIDs []int64) (propertyIDResourceID map[int64]int64, err error)
	GetPropertyIDsMap  func(resourceIDs []int64) (resourceIDPropertyID map[int64]int64, err error)

	ObjectLoader               func(tx isql.SqlTxStore, resourceIDs []int64, typ entity.ResourceType) (resources []entity.ResourceCreator, count int64, err error)
	GenerateGetPropertyIDMapFn func(tx isql.SqlTxStore) GetPropertyIDsMap

	PropertyCheckExistMethod string
)

const (
	PropertyCheckExistMethodUpdate PropertyCheckExistMethod = "update"
	PropertyCheckExistMethodCreate PropertyCheckExistMethod = "create"
)

func (a *resource) createResource(tx isql.SqlTxStore, typ entity.ResourceType) (resourceID int64, err error) {
	return tx.Resource().Create(typ)
}

func (a *resource) CreateResource(parentID int64, parentTyp entity.ResourceType, typ entity.ResourceType, propertyCreator PropertyCreator, existCheckFn PropertyCheckExist) (resourceID int64, err error) {
	var finalResourceID int64 = 0

	err = a.store.SqlStore.Transaction(func(tx isql.SqlTxStore) (ie error) {
		// 检查是否存在
		if parentID > 0 {
			if parentTyp == 0 {
				return errors.WithCode(errors.CodeInvalidArguments, " parent type must specific")
			}
		}

		ie = existCheckFn(tx, parentID > 0, 0, parentID, parentTyp)
		if ie != nil {
			return ie
		}

		// 创建资源
		rid, ie := a.createResource(tx, typ)
		if ie != nil {
			return errors.Wrap(ie, errors.CodeUnexpect)
		}
		finalResourceID = rid
		// 创建属性+关联关系
		_, ie = propertyCreator(tx, rid)
		if ie != nil {
			return errors.Wrap(ie, errors.CodeUnexpect)
		}
		// 创建资源-资源关系
		if parentID > 0 {
			// 判断父级资源是否存在
			parent, err := tx.Resource().GetByID(parentID)
			if err != nil {
				return errors.Wrap(err, errors.CodeUnexpect)
			}
			if parent.ID <= 0 {
				return errors.WithCode(errors.CodeObjectNotFound, fmt.Sprintf("parent Resource = %d not exist", parentID))
			}
			// 检查父节点类型是否是我们指定的类型
			if parent.Type != parentTyp {
				return errors.WithCode(errors.CodeInvalidArguments, "resource type is invalid")
			}

			// 关联关系
			_, err = tx.ResourceResourceRel().Create(parentID, finalResourceID, entity.ResourceRelType(parentTyp.String()+"-"+typ.String()))
			if err != nil {
				return errors.Wrap(err, errors.CodeUnexpect)
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	return finalResourceID, nil
}

func (a *resource) DeleteResource(resourceID int64, deleter PropertyDeleter) (err error) {
	if resourceID <= 0 {
		return errors.WithCode(errors.CodeInvalidArguments, "resource ID must specific")
	}

	err = a.store.SqlStore.Transaction(func(tx isql.SqlTxStore) (ie error) {
		// 删除属性表
		ie = deleter(tx, resourceID)
		if ie != nil {
			return errors.Wrap(ie, errors.CodeUnexpect)
		}
		// 删除资源表
		return tx.Resource().DeleteByID(resourceID)
	})
	return err
}

func (a *resource) UpdateResource(resourceID int64, parentType entity.ResourceType, updater PropertyUpdater, existCheck PropertyCheckExist) (updated entity.ResourceCreator, err error) {
	err = a.store.SqlStore.Transaction(func(tx isql.SqlTxStore) (ie error) {
		// 检查是否存在
		if ie := existCheck(tx, parentType != 0, resourceID, 0, parentType); ie != nil {
			return errors.Wrap(err, errors.CodeUnexpect)
		}

		// 更新属性表
		updated, ie = updater(tx, resourceID)
		if ie != nil {
			return errors.Wrap(ie, errors.CodeUnexpect)
		}
		// 还原成资源ID
		updated.SetID(resourceID)
		return nil
	})
	return updated, err
}

func (a *resource) GetResource(resourceID int64, getter PropertyGetter) (info entity.ResourceCreator, err error) {
	if resourceID <= 0 {
		return nil, errors.WithCode(errors.CodeInvalidArguments, "resource ID must specific")
	}
	err = a.store.SqlStore.Transaction(func(tx isql.SqlTxStore) (ie error) {
		info, ie = getter(tx, resourceID)
		if ie != nil {
			return errors.Wrap(ie, errors.CodeUnexpect)
		}
		// 还原成资源ID
		info.SetID(resourceID)
		return nil
	})
	return info, err
}

func (a *resource) ListResources(parentID int64, parentTyp entity.ResourceType, childTyp entity.ResourceType, filter entity.Filter, lister PropertyLister) (resources []entity.ResourceCreator, count int64, err error) {
	err = a.store.SqlStore.Transaction(func(tx isql.SqlTxStore) (ie error) {
		resources, count, ie = a.listResources(tx, parentID, parentTyp, childTyp, filter, lister)
		if ie != nil {
			return errors.Wrap(err, errors.CodeUnexpect)
		}
		return nil
	})
	return resources, count, err
}

func (a *resource) listResources(tx isql.SqlTxStore, parentID int64, parentTyp entity.ResourceType, childTyp entity.ResourceType, filter entity.Filter, lister PropertyLister) (resources []entity.ResourceCreator, count int64, err error) {
	// 查询直系子节点
	if parentID > 0 {
		children, err := tx.ResourceResourceRel().GetDirectChildren(parentID, entity.ResourceRelType(parentTyp.String()+"-"+childTyp.String()))
		if err != nil {
			return nil, 0, errors.Wrap(err, errors.CodeUnexpect)
		}
		childResourceIDs := []int64{}
		for _, v := range children {
			childResourceIDs = append(childResourceIDs, v.ID)
		}
		// 去重
		childResourceIDs = tools.Int64SliceDuplicate(childResourceIDs)

		if len(childResourceIDs) > 0 {
			filter = filter.And("ResourceID", entity.In, childResourceIDs)
		} else {
			// 没有子节点，则不用继续查询
			return nil, 0, errors.WithCode(errors.CodeObjectNotFound, fmt.Sprintf(`rsource=%d no children resource`, parentID))
		}
	}

	var getResourceMapFn GetResourceIDsMap
	resources, count, getResourceMapFn, err = lister(tx, filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeUnexpect)
	}
	// 转换ID
	propertyIDs := make([]int64, 0, len(resources))
	for _, v := range resources {
		propertyIDs = append(propertyIDs, v.GetPropertyID())
	}
	// 去重
	propertyIDs = tools.Int64SliceDuplicate(propertyIDs)

	propertyIDResourceID, ie := getResourceMapFn(propertyIDs)
	if ie != nil {
		return nil, 0, errors.Wrap(ie, errors.CodeUnexpect)
	}
	for i, v := range resources {
		resourceID := propertyIDResourceID[v.GetPropertyID()]
		v.SetID(resourceID)
		resources[i] = v
	}
	return resources, count, err
}

func (a *resource) ResourceLoader(parentID int64, parentType entity.ResourceType, filter entity.Filter, genGetPropertyIDsMapFn GenerateGetPropertyIDMapFn, lister PropertyLister, notIn bool) ObjectLoader {
	return func(tx isql.SqlTxStore, resourceIDs []int64, typ entity.ResourceType) (resources []entity.ResourceCreator, count int64, err error) {
		// 资源ID转化为属性ID
		resourcePropertyIDMap, err := genGetPropertyIDsMapFn(tx)(resourceIDs)
		if err != nil {
			return nil, 0, errors.Wrap(err, errors.CodeUnexpect)
		}
		propertyIDs := make([]int64, 0, len(resourcePropertyIDMap))
		for _, v := range resourcePropertyIDMap {
			propertyIDs = append(propertyIDs, v)
		}
		// 去重
		propertyIDs = tools.Int64SliceDuplicate(propertyIDs)

		if notIn {
			filter = filter.And("ID", entity.NotIn, propertyIDs)
		} else {

			filter = filter.And("ID", entity.In, propertyIDs)
		}
		return a.listResources(tx, parentID, parentType, typ, filter, lister)
	}
}
