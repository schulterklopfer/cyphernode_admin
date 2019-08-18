package queries

import (
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func GetUser( id uint, recursive bool ) (*models.UserModel, error) {
  db := dataSource.GetDB()
  var user models.UserModel
  db.Take(&user, id)

  if user.ID > 0 {
    if recursive {
      err := loadRolesForUser(&user)
      if err != nil {
        return nil, err
      }
    }
    return &user, nil
  } else {
    return nil, nil
  }

}

func FindUsers( where []interface{}, order string, limit int, offset uint, recursive bool ) ([]*models.UserModel, error) {

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
     db = db.Where( where[0].(string), where[1:] )
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

   if recursive {
     for i:=0; i<len(users); i++ {
       loadRolesForUser( users[i] )
     }
   }

   return users, db.Error

}


func loadRolesForUser( user *models.UserModel ) error {
  db := dataSource.GetDB()
  var roles []*models.RoleModel
  db.Model(&user).Association("Roles").Find(&roles)
  user.Roles = roles
  return db.Error
}