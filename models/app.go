package models

import (
  "database/sql/driver"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/jinzhu/gorm"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
)

type AccessPolicies []*storage.AccessPolicy

// this is used to create jsonb Data which then
// will ne saved to the db by gorm
func (aps AccessPolicies) Value() (driver.Value, error) {
  bytes, err := json.Marshal( aps )
  if err != nil {
    return nil, err
  }
  return bytes, nil
}

// convert jsonb Data from database back into
// a struct
func (aps *AccessPolicies) Scan(value interface{}) error {
  bytes, ok := value.([]byte)
  if !ok {
    return errors.New(fmt.Sprint("Failed to unmarshal access policies:", value))
  }
  err := json.Unmarshal(bytes, aps)
  return err
}

type AppModel struct {
  gorm.Model
  Hash           string         `json:"hash" gorm:"type:varchar(32);unique_index;not null"`
  Secret         string         `json:"-" gorm:"type:varchar(32);unique_index;not null"`
  MountPoint     string         `json:"mountPoint" gorm:"type:varchar(32);unique_index;not null"`
  Name           string         `json:"name" gorm:"type:varchar(30);not null" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_\\- ]+$"`
  Description    string         `json:"description" gorm:"type:varchar(255)"`
  AvailableRoles []*RoleModel   `json:"availableRoles" gorm:"foreignkey:AppId;preload"`
  AccessPolicies AccessPolicies `json:"accessPolicies,omitempty" gorm:"type:jsonb"`
}

func ( app *AppModel ) AfterDelete( tx *gorm.DB ) {
  var roles []RoleModel
  tx.Model(app).Association("AvailableRoles" ).Find(&roles)
  for i:=0; i< len(roles); i++ {
    tx.Delete( roles[i] )
    // Why do I have to call this manually?
    roles[i].AfterDelete( tx )
  }
}

func ( app *AppModel ) BeforeDelete( tx *gorm.DB ) (err error) {
  // very important. if no check, will delete all users if ID == 0
  if app.ID == 0 {
    err = cnaErrors.ErrNoSuchApp
    return
  }
  return
}