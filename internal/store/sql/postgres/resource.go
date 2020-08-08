package postgres

import (
	"github.com/jinzhu/gorm"
	sentity "github.com/zouxinjiang/axes/internal/store/sql/entity"
	"github.com/zouxinjiang/axes/pkg/errors"
)

type (
	resource struct {
		errorDealer
		db    *gorm.DB
		table string
	}
	resourceMdl struct {
		ID   int64                 `gorm:"column:id"`
		Type *sentity.ResourceType `gorm:"column:type"`
	}
)

func (resourceMdl) TableName() string {
	return TableResource
}

func (m resourceMdl) toEntity() sentity.Resource {
	return sentity.Resource{
		ID:   m.ID,
		Type: *m.Type,
	}
}

func newResource(db *gorm.DB) *resource {
	return &resource{
		db:    db,
		table: TableResource,
	}
}

func (r resource) Create(typ sentity.ResourceType) (resourceID int64, err error) {
	d := resourceMdl{
		Type: &typ,
	}
	err = r.db.Table(r.table).Create(&d).Take(&d).Error
	err = r.deal(err)
	if err != nil {
		return 0, err
	}
	return d.ID, nil
}

func (r resource) DeleteByID(resourceID int64) error {
	if resourceID <= 0 {
		return errors.WithCode(errors.CodeInvalidArguments, "resource ID must specific")
	}
	err := r.db.Table(r.table).Where(resourceMdl{
		ID: resourceID,
	}).Delete(&resource{}).Error
	err = r.deal(err)
	if err != nil {
		return err
	}
	return nil
}

func (r resource) GetByID(resourceID int64) (item sentity.Resource, err error) {
	tmp := resourceMdl{}
	err = r.db.Table(r.table).Where(resourceMdl{
		ID: resourceID,
	}).Take(&tmp).Error
	err = r.deal(err)
	if err != nil {
		return sentity.Resource{}, err
	}
	return tmp.toEntity(), nil
}

func (r resource) GetByType(typ sentity.ResourceType) (list []sentity.Resource, err error) {
	items := []resourceMdl{}
	err = r.db.Table(r.table).Where(resourceMdl{
		Type: &typ,
	}).Find(&items).Error
	err = r.deal(err)
	if err != nil {
		return nil, err
	}
	list = make([]sentity.Resource, 0, len(items))
	for _, v := range items {
		list = append(list, v.toEntity())
	}

	return list, nil
}

func (r resource) UpdateByID(resourceID int64, typ *sentity.ResourceType) (updated sentity.Resource, err error) {
	if resourceID <= 0 {
		return sentity.Resource{}, errors.WithCode(errors.CodeInvalidArguments, " resource ID must specific")
	}

	item := resourceMdl{}
	err = r.db.Table(r.table).Where(resourceMdl{
		ID: resourceID,
	}).Updates(resourceMdl{
		Type: typ,
	}).Take(&item).Error
	err = r.deal(err)
	if err != nil {
		return sentity.Resource{}, err
	}
	return item.toEntity(), nil
}

type (
	resourceResourceRel struct {
		errorDealer
		db    *gorm.DB
		table string
	}
	resourceResourceRelMdl struct {
		ID       int64                   `gorm:"column:id;auto_increment"`
		Type     sentity.ResourceRelType `gorm:"column:type"`
		ParentID int64                   `gorm:"column:parent_id"`
		ChildID  int64                   `gorm:"column:child_id"`
	}
)

func (resourceResourceRelMdl) TableName() string {
	return TableResourceResourceRel
}

func (m resourceResourceRelMdl) toEntity() sentity.ResourceResourceRel {
	return sentity.ResourceResourceRel{
		ID:       m.ID,
		ParentID: m.ParentID,
		ChildID:  m.ChildID,
		RelType:  m.Type,
	}
}

func newResourceResourceRel(db *gorm.DB) *resourceResourceRel {
	return &resourceResourceRel{
		db:    db,
		table: TableResourceResourceRel,
	}
}

func (r resourceResourceRel) GetDirectParents(nodeID int64, typ sentity.ResourceRelType) (parents []sentity.Resource, err error) {
	res := []resourceMdl{}
	err = r.db.Table(TableResourceResourceRel+` AS rel`).
		Joins(`JOIN `+TableResource+` AS resource ON rel.parent_id = resource.id`).
		Where(`rel."type" = ? AND rel.child_id = ?`, typ, nodeID).
		Select(`resource.*`).Find(&res).Error
	err = r.deal(err)
	if err != nil {
		return nil, err
	}
	parents = make([]sentity.Resource, 0, len(res))
	for _, v := range res {
		parents = append(parents, v.toEntity())
	}

	return parents, nil
}

func (r resourceResourceRel) GetDirectChildren(nodeID int64, typ sentity.ResourceRelType) (children []sentity.Resource, err error) {
	res := []resourceMdl{}
	err = r.db.Table(TableResourceResourceRel+` AS rel`).
		Joins(`JOIN `+TableResource+` AS resource ON rel.child_id = resource.id`).
		Where(`rel."type" = ? AND rel.parent_id = ?`, typ, nodeID).
		Select(`resource.*`).Find(&res).Error
	err = r.deal(err)
	if err != nil {
		return nil, err
	}
	children = make([]sentity.Resource, 0, len(res))
	for _, v := range res {
		children = append(children, v.toEntity())
	}

	return children, nil
}

func (r resourceResourceRel) GetAllParents(nodeID int64, typ sentity.ResourceRelType) (parents []sentity.Resource, err error) {
	items := []resourceMdl{}
	sql := `
WITH RECURSIVE t AS (
	SELECT * FROM ` + r.table + ` AS rel WHERE rel.child_id = ? AND rel."type" = ?
	UNION ALL
	SELECT rel2.* FROM ` + r.table + ` AS rel2,t	WHERE rel2.child_id=t.parent_id AND rel2."type" = ? 
)
SELECT  resource.* FROM t JOIN ` + TableResource + ` AS resource ON resource.id = t.parent_id;`
	err = r.db.Raw(sql, nodeID, typ, typ).Find(&items).Error
	err = r.deal(err)
	if err != nil {
		return nil, err
	}
	parents = make([]sentity.Resource, 0, len(items))
	for _, v := range items {
		parents = append(parents, v.toEntity())
	}

	return parents, nil
}

func (r resourceResourceRel) GetAllChildren(nodeID int64, typ sentity.ResourceRelType) (children []sentity.Resource, err error) {
	items := []resourceMdl{}
	sql := `
WITH RECURSIVE t AS (
	SELECT * FROM ` + r.table + ` AS rel WHERE rel.parent_id = ? AND rel."type" = ?
	UNION ALL
	SELECT rel2.* FROM ` + r.table + ` AS rel2,t	WHERE rel2.parent_id=t.child_id AND AND rel2."type" = ? 
)
SELECT resource.* FROM t JOIN ` + TableResource + ` AS resource ON resource.id = t.parent_id;`
	err = r.db.Raw(sql, nodeID, typ, typ).Find(&items).Error
	err = r.deal(err)
	if err != nil {
		return nil, err
	}
	children = make([]sentity.Resource, 0, len(items))
	for _, v := range items {
		children = append(children, v.toEntity())
	}
	return children, nil
}

func (r resourceResourceRel) Create(parentID, childID int64, typ sentity.ResourceRelType) (relID int64, err error) {
	d := resourceResourceRelMdl{
		Type:     typ,
		ParentID: parentID,
		ChildID:  childID,
	}
	err = r.db.Table(r.table).Create(&d).Error
	err = r.deal(err)
	if err != nil {
		return 0, err
	}
	return d.ID, nil
}

func (r resourceResourceRel) DeleteByID(relID int64) error {
	if relID <= 0 {
		return errors.WithCode(errors.CodeInvalidArguments, "rel ID must specific")
	}
	err := r.db.Table(r.table).Where(resourceResourceRelMdl{
		ID: relID,
	}).Delete(resourceResourceRelMdl{}).Error
	err = r.deal(err)
	if err != nil {
		return err
	}
	return nil
}

func (r resourceResourceRel) DeleteByParentID(parentID int64, typ sentity.ResourceRelType) (rowAffect int64, err error) {
	if parentID <= 0 {
		return 0, errors.WithCode(errors.CodeInvalidArguments, "parent ID must specific")
	}
	db := r.db.Table(r.table).Where(resourceResourceRelMdl{
		Type:     typ,
		ParentID: parentID,
	}).Delete(resourceResourceRel{})
	err = db.Error
	rowAffect = db.RowsAffected
	err = r.deal(err)
	if err != nil {
		return rowAffect, err
	}
	return rowAffect, nil
}

func (r resourceResourceRel) DeleteByChildID(childID int64, typ sentity.ResourceRelType) (rowAffect int64, err error) {
	if childID <= 0 {
		return 0, errors.WithCode(errors.CodeInvalidArguments, "child ID must specific")
	}
	db := r.db.Table(r.table).Where(resourceResourceRelMdl{
		Type:    typ,
		ChildID: childID,
	}).Delete(resourceResourceRelMdl{})
	err = db.Error
	rowAffect = db.RowsAffected
	err = r.deal(err)
	if err != nil {
		return rowAffect, err
	}
	return rowAffect, nil
}

func (r resourceResourceRel) DeleteByParentIDAndChildID(parentID, childID int64, typ sentity.ResourceRelType) (rowAffect int64, err error) {
	if childID <= 0 {
		return 0, errors.WithCode(errors.CodeInvalidArguments, "child ID must specific")
	}
	if parentID <= 0 {
		return 0, errors.WithCode(errors.CodeInvalidArguments, "parent ID must specific")
	}
	db := r.db.Table(r.table).Where(resourceResourceRelMdl{
		Type:     typ,
		ChildID:  childID,
		ParentID: parentID,
	}).Delete(resourceResourceRelMdl{})
	err = db.Error
	rowAffect = db.RowsAffected
	err = r.deal(err)
	if err != nil {
		return rowAffect, err
	}
	return rowAffect, nil
}

func (r resourceResourceRel) GetByID(ID int64) (e sentity.ResourceResourceRel, err error) {
	item := resourceResourceRelMdl{}
	err = r.db.Table(r.table).Where(resourceResourceRelMdl{
		ID: ID,
	}).Take(&item).Error
	err = r.deal(err)
	if err != nil {
		return sentity.ResourceResourceRel{}, err
	}
	return item.toEntity(), nil
}

func (r resourceResourceRel) GetByType(typ sentity.ResourceRelType) (list []sentity.ResourceResourceRel, err error) {
	items := []resourceResourceRelMdl{}
	err = r.db.Table(r.table).Where(resourceResourceRelMdl{
		Type: typ,
	}).Find(&items).Error
	err = r.deal(err)
	if err != nil {
		return nil, err
	}
	list = make([]sentity.ResourceResourceRel, 0, len(items))
	for _, v := range items {
		list = append(list, v.toEntity())
	}

	return list, nil
}

func (r resourceResourceRel) GetByParentIDAndChildID(parentID, childID int64, typ sentity.ResourceRelType) (e sentity.ResourceResourceRel, err error) {
	item := resourceResourceRelMdl{}
	err = r.db.Table(r.table).Where(resourceResourceRelMdl{
		Type:     typ,
		ParentID: parentID,
		ChildID:  childID,
	}).Take(&item).Error
	err = r.deal(err)
	if err != nil {
		return sentity.ResourceResourceRel{}, err
	}
	return item.toEntity(), nil
}

func (r resourceResourceRel) GetByParentID(parentID int64, typ sentity.ResourceRelType) (list []sentity.ResourceResourceRel, err error) {
	if parentID <= 0 {
		return nil, errors.WithCode(errors.CodeInvalidArguments, "parent ID must specific")
	}
	items := []resourceResourceRelMdl{}
	err = r.db.Table(r.table).Where(resourceResourceRelMdl{
		Type:     typ,
		ParentID: parentID,
	}).Find(&items).Error
	err = r.deal(err)
	if err != nil {
		return nil, err
	}
	list = make([]sentity.ResourceResourceRel, 0, len(items))
	for _, v := range items {
		list = append(list, v.toEntity())
	}
	return list, nil
}

func (r resourceResourceRel) GetByChildID(childID int64, typ sentity.ResourceRelType) (list []sentity.ResourceResourceRel, err error) {
	if childID <= 0 {
		return nil, errors.WithCode(errors.CodeInvalidArguments, "child ID must specific")
	}
	items := []resourceResourceRelMdl{}
	err = r.db.Table(r.table).Where(resourceResourceRelMdl{
		Type:    typ,
		ChildID: childID,
	}).Find(&items).Error
	err = r.deal(err)
	if err != nil {
		return nil, err
	}
	list = make([]sentity.ResourceResourceRel, 0, len(items))
	for _, v := range items {
		list = append(list, v.toEntity())
	}
	return list, nil
}
