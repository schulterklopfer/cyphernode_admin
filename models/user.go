package models

import (
  "errors"
  "github.com/jinzhu/gorm"
)

type UserModel struct {
  gorm.Model
  Login string `json:"login" gorm:"type:varchar(30);unique_index;not null" form:"login" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_\\-]+$"`
  Name string `json:"name" form:"name"` // optional
  Password string `json:"password" gorm:"type:varchar(128);not null" form:"password" validate:"nonzero" sbjt:"hashPassword"`
  EmailAddress string `json:"email_address" gorm:"type:varchar(100)" form:"emailAddress" validate:"max=100,regexp=(^$|^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"`
  Roles []*RoleModel `json:"roles" gorm:"many2many:user_roles;association_autoupdate:false" form:"roles" validate:"-"`
}

var ErrDuplicateUser = errors.New("user already exists")
var ErrUserHasUnknownRole = errors.New("user has unknown role")
var ErrNoSuchUser = errors.New( "no such user" )

func (user *UserModel) AfterCreate( tx *gorm.DB ) (err error) {
  var allAutoAssignRoles []*RoleModel
  tx.Where( &RoleModel{ AutoAssign: true }).Find( &allAutoAssignRoles )
  for i:=0; i< len(allAutoAssignRoles); i++ {
    tx.Model(user).Association("Roles").Append(allAutoAssignRoles[i])
  }
  return
}

func (user *UserModel) BeforeDelete( tx *gorm.DB ) (err error) {
  // very important. if no check, will delete all users if ID == 0
  if user.ID == 0 {
    err = ErrNoSuchUser
    return
  }
  return
}

func (user *UserModel) AfterDelete( tx *gorm.DB ) (err error) {
  tx.Model(user).Association("Roles").Clear()
  return
}

func (user *UserModel) BeforeSave( tx *gorm.DB ) (err error) {
  err = user.checkDuplicate(tx)
  if err != nil {
    return
  }
  err = user.checkRoles(tx)
  if err != nil {
    return
  }
  return
}

func (user *UserModel) BeforeCreate( tx *gorm.DB ) (err error) {
  err = user.checkDuplicate(tx)
  if err != nil {
    return
  }
  err = user.checkRoles(tx)
  if err != nil {
    return
  }
  return
}

func (user *UserModel) checkDuplicate( tx *gorm.DB ) error {
  var existingUsers []UserModel
  tx.Limit(1).Find( &existingUsers, "login = ? AND id != ?", user.Login, user.ID )

  if len(existingUsers) > 0 {
    return ErrDuplicateUser
  }
  return nil
}

func (user *UserModel) checkRoles( tx *gorm.DB ) error {
  for i:=0; i<len(user.Roles ); i++ {
    if user.Roles[i].ID == 0 {
      return ErrUserHasUnknownRole
    }
    var role RoleModel
    tx.Take( &role,  user.Roles[i].ID )
    if role.ID != user.Roles[i].ID {
      return ErrUserHasUnknownRole
    }
  }
  return nil
}