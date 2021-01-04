package cyphernodeAdmin

import (
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
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

  hashedPassword, err := password.HashPassword( cyphernodeAdmin.Config.InitialAdminPassword )
  if err != nil {
    return err
  }

  tx := db.Begin()

  if adminRole.ID != 1 {
    logwrapper.Logger().Info("adding admin role")
    adminRole.ID = 1
    adminRole.Name = ADMIN_APP_ADMIN_ROLE_NAME
    adminRole.Description = ADMIN_APP_ADMIN_ROLE_DESCRIPTION
    adminRole.AutoAssign = false
    adminRole.AppId = 1
    tx.Create(adminRole)
  }

  if adminApp.ID != 1 {
    logwrapper.Logger().Info("adding admin app")
    adminApp.ID = 1
    adminApp.Name = ADMIN_APP_NAME
    adminApp.Description = ADMIN_APP_DESCRIPTION
    adminApp.Hash = buildAdminHash( adminApp.Name, globals.CYPHERAPPS_REPO )
    adminApp.MountPoint = globals.BASE_ADMIN_MOUNTPOINT
    adminApp.AccessPolicies = models.AccessPolicies{
      /* General stuff */
      {
        Patterns: []string{".+"},
        Roles: []string{"*"},
        Actions: []string{"options"},
        Effect: "allow",
      },
      {
        Patterns: []string{"favicon.ico$"},
        Roles: []string{"*"},
        Actions: []string{"get", "options"},
        Effect: "allow",
      },
      /* API endpoints */
      {
        Patterns: []string{"^\\/api\\/v0\\/login$"},
        Roles: []string{"*"},
        Actions: []string{"post"},
        Effect: "allow",
      },
      {
        Patterns: []string{"^\\/api\\/v0\\/users"},
        Roles: []string{"admin"},
        Actions: []string{"get","post","patch","delete"},
        Effect: "allow",
      },
      {
        Patterns: []string{"^\\/api\\/v0\\/apps","^\\/api\\/v0\\/status"},
        Roles: []string{"*"},
        Actions: []string{"get"},
        Effect: "allow",
      },
      {
        Patterns: []string{"^\\/api\\/v0\\/apps"},
        Roles: []string{"admin"},
        Actions: []string{"post","patch"},
        Effect: "allow",
      },
      {
        Patterns: []string{"^\\/$", "^\\/_\\/"},
        Roles: []string{"*"},
        Actions: []string{"get"},
        Effect: "allow",
      },
      {
        Patterns: []string{"^\\/_\\/"},
        Roles: []string{"user","admin"},
        Actions: []string{"get"},
        Effect: "allow",
      },
    }
    if adminApp.Hash == "" {
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
    logwrapper.Logger().Info("adding admin user")
    adminUser.ID = 1
    adminUser.Login = cyphernodeAdmin.Config.InitialAdminLogin
    adminUser.Password = hashedPassword
    adminUser.Name = cyphernodeAdmin.Config.InitialAdminName
    adminUser.EmailAddress = cyphernodeAdmin.Config.InitialAdminEmailAddress
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

func buildAdminHash( label string, sourceLocation string ) string {
  bytes := make( []byte, 0 )
  bytes = append( bytes, []byte(label)... )
  bytes = append( bytes, []byte(sourceLocation)... )
  return helpers.TrimmedRipemd160Hash( &bytes )
}

