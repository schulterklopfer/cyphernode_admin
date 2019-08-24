package queries

import (
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func Get( model interface{}, id uint, recursive bool ) error {
  db := dataSource.GetDB()
  db.Take(model, id)
  if recursive {
    err := loadRoles(model)
    if err != nil {
      return err
    }
  }
  return nil
}

func Find( out interface{}, where []interface{}, order string, limit int, offset uint, recursive bool ) error {

  /*
     where == nil -> no where
     order == "" -> no order
     limit == -1 -> no limit
     offset == 0 -> no offset
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

  db.Find( out )

  if recursive {
    switch out.(type) {
    case *[]*models.UserModel:
      users := *out.(*[]*models.UserModel)
      for i:=0; i<len(users); i++ {
        _ = loadRoles(users[i])
      }
    case *[]*models.AppModel:
      apps := *out.(*[]*models.AppModel)
      for i:=0; i<len(apps); i++ {
        _ = loadRoles(apps[i])
      }
    }
  }

  return db.Error

}

func loadRoles( in interface{} ) error {
  db := dataSource.GetDB()
  var roles []*models.RoleModel
  switch in.(type) {
  case *models.UserModel:
    if in.(*models.UserModel).ID > 0 {
      db.Model(in).Association("Roles").Find(&roles)
      in.(*models.UserModel).Roles = roles
    }
  case *models.AppModel:
    if in.(*models.AppModel).ID > 0 {
      db.Model(in).Association("AvailableRoles").Find(&roles)
      in.(*models.AppModel).AvailableRoles = roles
    }
  }
  return db.Error
}