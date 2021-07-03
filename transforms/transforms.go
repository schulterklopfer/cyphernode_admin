/*
 * MIT License
 *
 * Copyright (c) 2021 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package transforms

import (
  "github.com/SatoshiPortal/cam/storage"
  "github.com/schulterklopfer/cyphernode_fauth/models"
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
      out.(*AppV0).MountPoint = in.(*models.AppModel).MountPoint
      out.(*AppV0).Name = in.(*models.AppModel).Name
      out.(*AppV0).Description = in.(*models.AppModel).Description
      out.(*AppV0).Version = in.(*models.AppModel).Version
      out.(*AppV0).Meta = in.(*models.AppModel).Meta

      roleCount := len(in.(*models.AppModel).AvailableRoles)
      transformedAvailableRoles :=make( []*RoleV0, roleCount )

      for i:=0; i<roleCount; i++ {
        transformedAvailableRoles[i] = new(RoleV0)
        Transform(in.(*models.AppModel).AvailableRoles[i], transformedAvailableRoles[i])
      }
      out.(*AppV0).AvailableRoles = transformedAvailableRoles

      accessPolicyCount := len(in.(*models.AppModel).AccessPolicies)
      transformedAccessPolicies :=make( []*AccessPolicyV0, accessPolicyCount )

      for i:=0; i<accessPolicyCount; i++ {
        transformedAccessPolicies[i] = new(AccessPolicyV0)
        Transform(in.(*models.AppModel).AccessPolicies[i], transformedAccessPolicies[i])
      }

      out.(*AppV0).AccessPolicies = transformedAccessPolicies
      transformed = true
    }
  case *RoleV0:
    switch in.(type) {
    case *models.RoleModel:
      out.(*RoleV0).ID = in.(*models.RoleModel).ID
      out.(*RoleV0).AppId = in.(*models.RoleModel).AppId
      out.(*RoleV0).Name = in.(*models.RoleModel).Name
      out.(*RoleV0).Description = in.(*models.RoleModel).Description
      out.(*RoleV0).AutoAssign = in.(*models.RoleModel).AutoAssign
      transformed = true
    }
  case *AccessPolicyV0:
    switch in.(type) {
    case *storage.AccessPolicy:
      out.(*AccessPolicyV0).Effect = in.(*storage.AccessPolicy).Effect
      out.(*AccessPolicyV0).Actions = in.(*storage.AccessPolicy).Actions
      out.(*AccessPolicyV0).Roles = in.(*storage.AccessPolicy).Roles
      out.(*AccessPolicyV0).Patterns = in.(*storage.AccessPolicy).Patterns
      transformed = true
    }
  }
  return transformed
}

