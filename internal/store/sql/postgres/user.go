package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/zouxinjiang/axes/internal/store/sql/entity"
	"time"
)

type (
	UserMdl struct {
		ID           int64     `gorm:"Column:id"`
		UserName     string    `gorm:"Column:user_name"`
		DisplayName  *string   `gorm:"Column:display_name"`
		Password     []byte    `gorm:"Column:password"`
		Mobile       *string   `gorm:"Column:mobile"`
		Email        *string   `gorm:"Column:email"`
		UserSourceID int64     `gorm:"Column:user_source_id"`
		OriginID     *string   `gorm:"Column:origin_id"`
		CreatedAt    time.Time `gorm:"Column:created_at"`
		UpdatedAt    time.Time `gorm:"Column:updated_at"`
	}
	User struct {
		errorDealer
		objectFuncTemplate
		db    *gorm.DB
		table string
	}
)

func (UserMdl) TableName() string {
	return TableUser
}

func (m UserMdl) ToEntity() entity.User {
	return entity.User{
		ID:           m.ID,
		UserName:     m.UserName,
		DisplayName:  m.DisplayName,
		Password:     m.Password,
		Mobile:       m.Mobile,
		Email:        m.Email,
		UserSourceID: m.UserSourceID,
		OriginID:     m.OriginID,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func newUser(db *gorm.DB) *User {
	return &User{
		db:    db,
		table: TableUser,
	}
}

func (u User) GetResourceIDMap(propertyIDs []int64) (propertyIDResourceIDMap map[int64]int64, err error) {
	propertyIDResourceIDMap, err = u.objectGetResourceIDMapTemplate(u.db, propertyIDs, paramsResourcePropertyRel{
		resourceTableName:               TableResource,
		resourcePropertyRelTableName:    TableResourceUserRel,
		propertyTableName:               TableUser,
		relTableAttachResourceFieldName: "resource_id",
		relTableAttachPropertyFieldName: "user_id",
	})
	err = u.deal(err)
	return propertyIDResourceIDMap, err
}

func (u User) GetPropertyIDMap(resourceIDs []int64) (resourceIDPropertyIDMap map[int64]int64, err error) {
	resourceIDPropertyIDMap, err = u.objectGetPropertyIDMapTemplate(u.db, resourceIDs, paramsResourcePropertyRel{
		resourceTableName:               TableResource,
		resourcePropertyRelTableName:    TableResourceUserRel,
		propertyTableName:               TableUser,
		relTableAttachResourceFieldName: "resource_id",
		relTableAttachPropertyFieldName: "user_id",
	})
	err = u.deal(err)
	return resourceIDPropertyIDMap, err
}

func (u User) Create(username, displayName string, password []byte, mobile, email string, sourceID int64, originID string) (userID int64, err error) {
	m := UserMdl{
		UserName:     username,
		DisplayName:  &displayName,
		Password:     password,
		Mobile:       &mobile,
		Email:        &email,
		UserSourceID: sourceID,
		OriginID:     &originID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = u.db.Table(u.table).Create(m).Take(&m).Error
	err = u.deal(err)
	return m.ID, err
}

func (u User) Update(userID int64, password []byte, displayName, mobile, email *string) (updated entity.User, err error) {
	m := UserMdl{}
	err = u.db.Table(u.table).Where(UserMdl{
		ID: userID,
	}).Updates(UserMdl{
		DisplayName: displayName,
		Password:    password,
		Mobile:      mobile,
		Email:       email,
		UpdatedAt:   time.Now(),
	}).Take(&m).Error
	err = u.deal(err)
	if err != nil {
		return entity.User{}, err
	}
	return m.ToEntity(), nil
}

func (u User) DeleteByID(userID int64) (err error) {
	panic("implement me")
}

func (u User) List(filter entity.Filter) (list []entity.User, cnt int64, err error) {
	panic("implement me")
}

func (u User) ListResource(filter entity.Filter) (list []entity.User, cnt int64, err error) {
	panic("implement me")
}

func (u User) GetByID(userID int64) (entity.User, error) {
	panic("implement me")
}
