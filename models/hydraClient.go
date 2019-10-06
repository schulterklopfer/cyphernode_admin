package models

import (
	"github.com/jinzhu/gorm"
)

type HydraClientModel struct {
	gorm.Model
	App *AppModel `json:"-" gorm:"foreignkey:AppID" form:"-" validate:"-"`
	AppID uint `json:"-" gorm:"DEFAULT:0" form:"-" validate:"-"`
	ClientID string `json:"-" gorm:"type:varchar(100)" form:"-" validate:"-"`
	Secret string `json:"-" gorm:"type:varchar(32)" form:"-" validate:"-"`
}
