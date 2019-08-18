package queries

import (
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func GetUser( id uint, recursive bool ) *models.UserModel {
  db := dataSource.GetDB()
  var user models.UserModel
  db.Take(&user, id)

  if user.ID > 0 {
    if recursive {
      loadRolesForUser(&user)
    }
    return &user
  } else {
    return nil
  }

}

func FindUsers( where []string, order string, limit int, offset uint ) []*models.UserModel {

  /*
    where == nil -> no where
    order == "" -> no order
    limit == -1 -> no limit
    offset == 0 -> no offset

    Example:

    users = queries.FindUsers( []string{"name LIKE ?", "name%"}, "name", -1,0)

  */

   db := dataSource.GetDB()

   if len(where) > 0 {
     db = db.Where( where[0], where[1:] )
   }

   if order != "" {
     db = db.Order( order )
   }

   if limit != -1 {
     db = db.Limit( limit )
   }

   if offset > 0 {
     db = db.Offset( offset )
   }

   var users []*models.UserModel

   db.Find( &users )

   return users

}


func loadRolesForUser( user *models.UserModel ) {
  db := dataSource.GetDB()
  var roles []*models.RoleModel
  db.Model(&user).Association("Roles").Find(&roles)
  user.Roles = roles
}