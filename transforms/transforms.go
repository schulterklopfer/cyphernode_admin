package transforms

import "github.com/schulterklopfer/cyphernode_admin/dataSource/models"

type (
  AppV0 struct {
    ID uint   `json:"id"`
    Hash string `json:"hash"`
    Name string `json:"name"`
    Description string `json:"description"`
    AvailableRoles []*RoleV0 `json:"availableRoles"`
  }

  RoleV0 struct {
    ID uint   `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    AutoAssign bool `json:"autoAssign"`
    AppId uint `json:"appId"`
  }

  UserV0 struct {
    ID uint   `json:"id"`
    Login string `json:"login"`
    Name string `json:"name"`
    EmailAddress string `json:"emailAddress"`
    Roles []*RoleV0 `json:"roles"`
  }
)

// transform from database model to api object
func Transform( in interface{}, out interface{} ) bool {
  transformed := false
  switch out.(type) {
  case *UserV0:
    switch in.(type) {
    case *models.UserModel:
      out.(*UserV0).ID = in.(*models.UserModel).ID
      out.(*UserV0).Login = in.(*models.UserModel).Login
      out.(*UserV0).Name = in.(*models.UserModel).Name
      out.(*UserV0).EmailAddress = in.(*models.UserModel).EmailAddress

      roleCount := len(in.(*models.UserModel).Roles)
      transformedRoles :=make( []*RoleV0, roleCount )

      for i:=0; i<roleCount; i++ {
        transformedRoles[i] = new(RoleV0)
        Transform(in.(*models.UserModel).Roles[i], transformedRoles[i])
      }

      out.(*UserV0).Roles = transformedRoles
      transformed = true
    }
  case *AppV0:
    switch in.(type) {
    case *models.AppModel:
      out.(*AppV0).ID = in.(*models.AppModel).ID
      out.(*AppV0).Hash = in.(*models.AppModel).Hash
      out.(*AppV0).Name = in.(*models.AppModel).Name
      out.(*AppV0).Description = in.(*models.AppModel).Description

      roleCount := len(in.(*models.AppModel).AvailableRoles)
      transformedAvailableRoles :=make( []*RoleV0, roleCount )

      for i:=0; i<roleCount; i++ {
        transformedAvailableRoles[i] = new(RoleV0)
        Transform(in.(*models.AppModel).AvailableRoles[i], transformedAvailableRoles[i])
      }

      out.(*AppV0).AvailableRoles = transformedAvailableRoles
      transformed = true
    }
  case *RoleV0:
    switch in.(type) {
    case *models.RoleModel:
      out.(*RoleV0).ID = in.(*models.RoleModel).ID
      out.(*RoleV0).Name = in.(*models.RoleModel).Name
      out.(*RoleV0).Description = in.(*models.RoleModel).Description
      out.(*RoleV0).AutoAssign = in.(*models.RoleModel).AutoAssign
      out.(*RoleV0).AppId = in.(*models.RoleModel).AppId
      transformed = true
    }
  }
  return transformed
}

