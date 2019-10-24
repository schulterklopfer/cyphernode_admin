package cyphernodeAdmin

import (
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/password"
  "github.com/schulterklopfer/cyphernode_admin/queries"
)

func (cyphernodeAdmin *CyphernodeAdmin) migrate() error {

  // Create adminUser id=1, cyphernodeAdmin id=1, adminRole id=1
  adminRole := new(models.RoleModel)
  adminApp := new(models.AppModel)
  adminUser := new(models.UserModel)

  db := dataSource.GetDB()

  _ = queries.Get( adminRole, 1, true )
  _ = queries.Get( adminApp, 1, true )
  _ = queries.Get( adminUser, 1, true )

  hashedPassword, err := password.HashPassword( cyphernodeAdmin.config.InitialAdminPassword )
  if err != nil {
    return err
  }

  roles := make( []*models.RoleModel, 1 )

  tx := db.Begin()

  if adminRole.ID != 1 {
    adminRole.ID = 1
    adminRole.Name = ADMIN_APP_ADMIN_ROLE_NAME
    adminRole.Description = ADMIN_APP_ADMIN_ROLE_DESCRIPTION
    adminRole.AutoAssign = false
    adminRole.AppId = 1
    tx.Create(adminRole)
  }

  roles[0]= adminRole

  if adminApp.ID != 1 {
    adminApp.ID = 1
    adminApp.Name = ADMIN_APP_NAME
    adminApp.Description = ADMIN_APP_DESCRIPTION
    adminApp.ClientSecret = helpers.RandomString(32 )
    adminApp.ClientID = helpers.RandomString(32 )
    if adminApp.ClientSecret == "" || adminApp.ClientID == "" {
      return cnaErrors.ErrMigrationFailed
    }
    tx.Create(adminApp)
  }

  hasAdminRole := false
  for i:=0; i<len(adminApp.AvailableRoles); i++ {
    if adminApp.AvailableRoles[i].ID == adminRole.ID {
      hasAdminRole = true
    }
  }

  if !hasAdminRole {
    tx.Model(&adminApp).Association("AvailableRoles").Append(adminRole)
  }

  if adminUser.ID != 1 {
    adminUser.ID = 1
    adminUser.Login = cyphernodeAdmin.config.InitialAdminLogin
    adminUser.Password = hashedPassword
    adminUser.Name = cyphernodeAdmin.config.InitialAdminName
    adminUser.EmailAddress = cyphernodeAdmin.config.InitialAdminEmailAddress
    tx.Create(adminUser)
  }

  hasAdminRole = false
  for i:=0; i<len(adminUser.Roles); i++ {
    if adminUser.Roles[i].ID == adminRole.ID {
      hasAdminRole = true
    }
  }

  if !hasAdminRole {
    tx.Model(&adminUser).Association("Roles").Append(adminRole)
  }

  return tx.Commit().Error

}

