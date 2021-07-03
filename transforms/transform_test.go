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

package transforms_test

import (
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "strconv"
  "testing"
)

func TestUserModelTransformV0(t *testing.T) {
  userModel := new( models.UserModel )
  var transformed transforms.UserV0

  userModel.ID = 1
  userModel.Login = "login"
  userModel.Name = "name"
  userModel.EmailAddress = "email address"

  roles := make( []*models.RoleModel,2)

  for i := 0; i<2; i++ {
    roles[i] = new(models.RoleModel)
    roles[i].Name = "role"+strconv.Itoa(i)
    roles[i].Description = "description"+strconv.Itoa(i)
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

func TestAppModelTransformV0(t *testing.T) {
  appModel := new( models.AppModel )
  var transformed transforms.AppV0

  appModel.ID = 1
  appModel.Name = "name"
  appModel.Description = "description"

  availableRoles := make( []*models.RoleModel,2)

  for i := 0; i<2; i++ {
    availableRoles[i] = new(models.RoleModel)
    availableRoles[i].Name = "role"+strconv.Itoa(i)
    availableRoles[i].Description = "description"+strconv.Itoa(i)
    availableRoles[i].AutoAssign = (i+1)%2==0
    availableRoles[i].AppId = 1
  }

  appModel.AvailableRoles = availableRoles

  transforms.Transform(appModel, &transformed )

  if appModel.ID != transformed.ID ||
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

func TestRoleModelTransformV0(t *testing.T) {
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