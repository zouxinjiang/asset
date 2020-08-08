package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/zouxinjiang/axes/internal/store/sql/entity"
	"github.com/zouxinjiang/axes/pkg/errors"
)

type (
	errorDealer struct {
	}

	// 对象函数模板
	objectFuncTemplate struct {
	}

	// 关系函数模板
	relationFuncTemplate struct {
	}

	paramsResourcePropertyRel struct {
		resourceTableName               string
		resourcePropertyRelTableName    string
		propertyTableName               string
		relTableAttachResourceFieldName string
		relTableAttachPropertyFieldName string
	}

	paramsListResource struct {
		resourceTableName               string
		resourcePropertyRelTableName    string
		propertyTableName               string
		relTableAttachResourceFieldName string
		relTableAttachPropertyFieldName string
		resourceType                    entity.ResourceType
	}

	paramsListResourceWithRRRelProperty struct {
		resourceTableName               string
		resourcePropertyRelTableName    string
		propertyTableName               string
		relTableAttachResourceFieldName string
		relTableAttachPropertyFieldName string
		resourceType                    entity.ResourceType

		// 关系表关联资源表的字段名
		RRRelTableAttachResourceFeildName string
		//关系表
		RRRelTableName string
		// 关系属性表
		RRRelPropertyTableName string
		// 关系属性表关联关系表字段名
		RRRelPropertyAttachRRRelFieldName string
	}

	paramsGetResourceTemplate struct {
		resourceTableName               string
		resourcePropertyRelTableName    string
		relTableAttachResourceFieldName string
		relTableAttachPropertyFieldName string
	}

	paramsGetPropertyIDsTemplate struct {
		resourcePropertyRelTableName    string
		propertyTableName               string
		relTableAttachResourceFieldName string
		relTableAttachPropertyFieldName string
	}
)

func (errorDealer) deal(err error) error {
	if gorm.IsRecordNotFoundError(err) {
		return errors.Wrap(err, errors.CodeObjectNotFound)
	}
	if err != nil {
		return errors.Wrap(err, errors.CodeUnexpect)
	}
	return err
}

// 属性ID转换成资源ID
func (objectFuncTemplate) objectGetResourceIDMapTemplate(db *gorm.DB, propertyIDs []int64, params paramsResourcePropertyRel) (propertyIDResourceIDMap map[int64]int64, err error) {
	if len(propertyIDs) <= 0 {
		return nil, nil
	}
	type mdl struct {
		ResourceID int64 `gorm:"column:resource_id"`
		PropertyID int64 `gorm:"column:property_id"`
	}

	items := []mdl{}
	err = db.Table(params.propertyTableName+` AS property`).
		Joins(`JOIN `+params.resourcePropertyRelTableName+` AS rel ON property.id = rel.`+params.relTableAttachPropertyFieldName).
		Joins(`JOIN `+params.resourceTableName+` AS resource ON resource.id = rel.`+params.relTableAttachResourceFieldName).
		Where(`property.id  IN (?)`, propertyIDs).
		Select(`property.id AS property_id, resource.id AS resource_id`).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	propertyIDResourceIDMap = map[int64]int64{}
	for _, v := range items {
		propertyIDResourceIDMap[v.PropertyID] = v.ResourceID
	}
	return propertyIDResourceIDMap, nil
}

// 资源ID转换成属性ID
func (objectFuncTemplate) objectGetPropertyIDMapTemplate(db *gorm.DB, resourceIDs []int64, params paramsResourcePropertyRel) (resourceIDPropertyIDMap map[int64]int64, err error) {
	if len(resourceIDs) <= 0 {
		return nil, nil
	}
	type mdl struct {
		ResourceID int64 `gorm:"column:resource_id"`
		PropertyID int64 `gorm:"column:property_id"`
	}

	items := []mdl{}
	err = db.Table(params.propertyTableName+` AS property`).
		Joins(`JOIN `+params.resourcePropertyRelTableName+` AS rel ON property.id = rel.`+params.relTableAttachPropertyFieldName).
		Joins(`JOIN `+params.resourceTableName+` AS resource ON resource.id = rel.`+params.relTableAttachResourceFieldName).
		Where(`resource.id  IN (?)`, resourceIDs).
		Select(`property.id AS property_id, resource.id AS resource_id`).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	resourceIDPropertyIDMap = map[int64]int64{}
	for _, v := range items {
		resourceIDPropertyIDMap[v.ResourceID] = v.PropertyID
	}
	return resourceIDPropertyIDMap, nil
}

// 读取对象接口
func (objectFuncTemplate) objectListTemplate(db *gorm.DB, payload interface{}, tableName string, filter entity.Filter, keymap map[string]string) (cnt int64, err error) {
	query := db.Table(tableName)
	where, vals := genSqlWithTable(tableName, filter, keymap)
	if where != "" {
		query = query.Where(where, vals...)
	}
	if err := query.Count(&cnt).Error; err != nil {
		return 0, errors.Wrap(err, errors.CodeUnexpect)
	}
	if cnt == 0 {
		return 0, nil
	}
	if filter.PageSize > 0 {
		query = query.Limit(filter.PageSize).Offset(filter.Offset())
	}
	query = WrapQueryOrder(query, filter, keymap)
	err = query.Find(&payload).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// 关联资源ID读取对象接口
func (objectFuncTemplate) objectListResourceTemplate(db *gorm.DB, payload interface{}, params paramsListResource, filter entity.Filter, keymap map[string]string) (cnt int64, err error) {
	query := db.Table(params.resourceTableName+` AS resource`).
		Joins(`JOIN `+params.resourcePropertyRelTableName+` AS rel ON rel.`+params.relTableAttachResourceFieldName+` = resource.id AND resource."type" = ?`, params.resourceType).
		Joins(`JOIN ` + params.propertyTableName + ` AS property ON rel.` + params.relTableAttachPropertyFieldName + `=property.id`)
	listKeyMap := map[string]string{}
	for k, v := range keymap {
		listKeyMap[k] = "property." + v
	}
	listKeyMap["ResourceID"] = "resource.id"

	where, vals := genSqlWithTable(params.propertyTableName, filter, listKeyMap)
	if where != "" {
		query = query.Where(where, vals...)
	}

	if err := query.Count(&cnt).Error; err != nil {
		return 0, errors.Wrap(err, errors.CodeUnexpect)
	}
	if cnt == 0 {
		return 0, nil
	}
	if filter.PageSize > 0 {
		query = query.Limit(filter.PageSize).Offset(filter.Offset())
	}
	query = WrapQueryOrder(query, filter, listKeyMap)

	err = query.Select(`resource.id AS resource_id,property.*`).Find(payload).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// 关联资源ID以及关系属性表读取对象接口
func (objectFuncTemplate) objectListResourceWithRRRelPropertyTemplate(db *gorm.DB, payload interface{}, params paramsListResourceWithRRRelProperty, filter entity.Filter, keymap map[string]string, RRRelPropertyKeyMap map[string]string) (cnt int64, err error) {
	query := db.Table(params.resourceTableName+` AS resource`).
		Joins(`JOIN `+params.resourcePropertyRelTableName+` AS rel ON rel.`+params.relTableAttachResourceFieldName+` = resource.id AND resource."type" = ?`, params.resourceType).
		Joins(`JOIN ` + params.propertyTableName + ` AS property ON rel.` + params.relTableAttachPropertyFieldName + `=property.id`).
		Joins(`LEFT JOIN ` + params.RRRelTableName + ` AS rrrel ON resource.id = rrrel.` + params.RRRelTableAttachResourceFeildName).
		Joins(`LEFT JOIN ` + params.RRRelPropertyTableName + ` AS relation ON rrrel.id = relation.` + params.RRRelPropertyAttachRRRelFieldName)

	listKeyMap := map[string]string{}
	for k, v := range keymap {
		listKeyMap[k] = "property." + v
	}

	for k, v := range RRRelPropertyKeyMap {
		listKeyMap[k] = "relation." + v
		listKeyMap["Relation."+k] = "relation." + v
	}

	listKeyMap["ResourceID"] = "resource.id"
	listKeyMap["Resource.ID"] = "resource.id"

	where, vals := genSqlWithTable(params.propertyTableName, filter, listKeyMap)
	if where != "" {
		query = query.Where(where, vals...)
	}

	if err := query.Count(&cnt).Error; err != nil {
		return 0, errors.Wrap(err, errors.CodeUnexpect)
	}
	if cnt == 0 {
		return 0, nil
	}
	if filter.PageSize > 0 {
		query = query.Limit(filter.PageSize).Offset(filter.Offset())
	}
	query = WrapQueryOrder(query, filter, listKeyMap)

	err = query.Select(`DISTINCT resource.id AS resource_id,property.*`).Find(payload).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// 属性ID读取资源
func (relationFuncTemplate) relationGetResourceTemplate(db *gorm.DB, propertyID int64, params paramsGetResourceTemplate) (list []entity.Resource, err error) {
	items := []resourceMdl{}
	err = db.Table(params.resourcePropertyRelTableName+` AS rel `).
		Joins(` LEFT JOIN `+params.resourceTableName+` AS resource ON rel.`+params.relTableAttachResourceFieldName+`=resource.id`).
		Select(` resource.*`).Where(`rel.`+params.relTableAttachPropertyFieldName+` = ?`, propertyID).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	list = []entity.Resource{}
	for _, v := range items {
		list = append(list, v.toEntity())
	}
	return list, nil
}

// 资源ID读取属性
func (relationFuncTemplate) relationGetPropertyTemplate(db *gorm.DB, resourceID int64, payload interface{}, params paramsGetPropertyIDsTemplate) error {
	err := db.Table(params.resourcePropertyRelTableName+` AS rel `).
		Joins(` LEFT JOIN `+params.propertyTableName+` AS property ON rel.`+params.relTableAttachPropertyFieldName+`=property.id`).
		Select(` property.*`).Where(`rel.`+params.relTableAttachResourceFieldName+` = ?`, resourceID).
		Find(payload).Error
	return err
}

// 资源ID读取属性ID
func (relationFuncTemplate) relationGetPropertyIDsTemplate(db *gorm.DB, resourceID int64, params paramsGetPropertyIDsTemplate) (list []int64, err error) {
	propertyIDs := []int64{}
	err = db.Table(params.resourcePropertyRelTableName+` AS rel `).
		Joins(` LEFT JOIN `+params.propertyTableName+` AS property ON rel.`+params.relTableAttachPropertyFieldName+`=property.id`).
		Where(`rel.`+params.relTableAttachResourceFieldName+` = ?`, resourceID).
		Pluck(`property.id`, &propertyIDs).Error
	if err != nil {
		return nil, err
	}

	return propertyIDs, nil
}
