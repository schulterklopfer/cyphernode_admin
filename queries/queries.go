package queries

import (
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func GetUser( id uint, recursive bool ) *models.UserModel {
  db := dataSource.GetDB()
  var user models.UserModel
  db.First(&user, id)

  if user.ID > 0 {
    if recursive {
      loadRolesForUser(&user)
    }
    return &user
  } else {
    return nil
  }

}

func loadRolesForUser( user *models.UserModel ) {
  db := dataSource.GetDB()
  var roles []*models.RoleModel
  db.Model(&user).Association("Roles").Find(&roles)
  user.Roles = roles
}