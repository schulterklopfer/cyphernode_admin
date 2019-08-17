package transforms_test

import (
  "github.com/schulterklopfer/cyphernode_admin/dataSource/models"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "testing"
)

func TestUserModelTransformV0( t* testing.T ) {
  userModel := new( models.UserModel )
  var transformed transforms.UserV0

  userModel.ID = 1
  userModel.Login = "login"
  userModel.Name = "name"
  userModel.EmailAddress = "email address"

  roles := make( []*models.RoleModel,2)

  for i := 0; i<2; i++ {
    roles[i] = new(models.RoleModel)
    roles[i].Name = "role"+string(i)
    roles[i].Description = "description"+string(i)
    roles[i].AutoAssign = (i+1)%2==0
    roles[i].AppId = 1
  }

  userModel.Roles = roles

  transforms.Transform( userModel, &transformed )

  if userModel.ID != transformed.ID ||
     userModel.Login != transformed.Login ||
     userModel.Name != transformed.Name ||
     userModel.EmailAddress != transformed.EmailAddress {
    t.Error("Property mismatch")
  }

  for i := 0; i<2; i++ {
    if transformed.Roles[i].Name != userModel.Roles[i].Name ||
       transformed.Roles[i].Description != userModel.Roles[i].Description ||
       transformed.Roles[i].AutoAssign != userModel.Roles[i].AutoAssign ||
       transformed.Roles[i].AppId != userModel.Roles[i].AppId {
      t.Error("Property mismatch")
    }
  }

}

func TestAppModelTransformV0( t* testing.T ) {
  appModel := new( models.AppModel )
  var transformed transforms.AppV0

  appModel.ID = 1
  appModel.Hash = "hash"
  appModel.Name = "name"
  appModel.Description = "description"

  availableRoles := make( []*models.RoleModel,2)

  for i := 0; i<2; i++ {
    availableRoles[i] = new(models.RoleModel)
    availableRoles[i].Name = "role"+string(i)
    availableRoles[i].Description = "description"+string(i)
    availableRoles[i].AutoAssign = (i+1)%2==0
    availableRoles[i].AppId = 1
  }

  appModel.AvailableRoles = availableRoles

  transforms.Transform(appModel, &transformed )

  if appModel.ID != transformed.ID ||
     appModel.Hash != transformed.Hash ||
     appModel.Name != transformed.Name ||
     appModel.Description != transformed.Description {
    t.Error("Property mismatch")
  }

  for i := 0; i<2; i++ {
    if transformed.AvailableRoles[i].Name != appModel.AvailableRoles[i].Name ||
       transformed.AvailableRoles[i].Description != appModel.AvailableRoles[i].Description ||
       transformed.AvailableRoles[i].AutoAssign != appModel.AvailableRoles[i].AutoAssign ||
       transformed.AvailableRoles[i].AppId != appModel.AvailableRoles[i].AppId {
      t.Error("Property mismatch")
    }
  }

}

func TestRoleModelTransformV0( t* testing.T ) {
  roleModel := new( models.RoleModel )
  var transformed transforms.RoleV0

  roleModel.ID = 1
  roleModel.Name = "name"
  roleModel.Description = "description"
  roleModel.AutoAssign = true
  roleModel.AppId = 1

  transforms.Transform(roleModel, &transformed )

  if roleModel.Name != transformed.Name ||
      roleModel.Description != transformed.Description ||
      roleModel.AutoAssign != transformed.AutoAssign ||
      roleModel.AppId != transformed.AppId {
    t.Error("Property mismatch")
  }
}