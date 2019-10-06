package cyphernodeAdmin

import (
	"github.com/schulterklopfer/cyphernode_admin/dataSource"
	"github.com/schulterklopfer/cyphernode_admin/models"
	"github.com/schulterklopfer/cyphernode_admin/password"
)

func (cyphernodeAdmin *CyphernodeAdmin) migrate() error {

	// Create adminUser id=1, cyphernodeAdmin id=1, adminRole id=1
	adminRole := new(models.RoleModel)
	adminApp := new(models.AppModel)
	adminUser := new(models.UserModel)

	db := dataSource.GetDB()

	db.Take(adminRole, 1 )
	db.Take(adminApp, 1 )
	db.Take(adminUser, 1 )

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
		adminApp.Hash = ADMIN_APP_HASH
		adminApp.AvailableRoles = roles
		tx.Create(adminApp)
	}

	if adminUser.ID != 1 {
		adminUser.ID = 1
		adminUser.Login = cyphernodeAdmin.config.InitialAdminLogin
		adminUser.Password = hashedPassword
		adminUser.Name = cyphernodeAdmin.config.InitialAdminName
		adminUser.EmailAddress = cyphernodeAdmin.config.InitialAdminEmailAddress
		adminUser.Roles = roles
		tx.Create(adminUser)
	}

	return tx.Commit().Error

}

