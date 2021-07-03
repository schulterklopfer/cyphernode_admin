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

import "github.com/schulterklopfer/cyphernode_fauth/models"

type (
  AppV0 struct {
    ID uint   `json:"id" form:"id"`
    Hash string `json:"hash" form:"clientid_like"`
    MountPoint string `json:"mountPoint" form:"mountpoint_like"`
    Name string `json:"name" form:"name_like"`
    Description string `json:"description" form:"description"`
    Version string `json:"version" form:"description"`
    AvailableRoles []*RoleV0 `json:"availableRoles" form:"availableRoles"`
    AccessPolicies []*AccessPolicyV0 `json:"accessPolicies" form:"availableRoles"`
    Meta *models.Meta `json:"meta,omitempty"`
  }

  RoleV0 struct {
    ID uint   `json:"id" form:"id"`
    AppId uint `json:"appId" form:"appId"`
    Name string `json:"name" form:"name"`
    Description string `json:"description" form:"description"`
    AutoAssign bool `json:"autoAssign" form:"autoAssign"`
  }

  AccessPolicyV0 struct {
    Roles []string `json:"roles"`
    Patterns []string `json:"resources"`
    Effect string `json:"effect"`
    Actions []string `json:"actions"`
  }

  // strange form names come from ng2_smart_table
  // didn't want to change the defaults
  UserV0 struct {
    ID uint   `json:"id" form:"id_like"`
    Login string `json:"login" form:"login_like"`
    Name string `json:"name" form:"name_like"`
    EmailAddress string `json:"email_address" form:"email_address_like"`
    Roles []*RoleV0 `json:"roles" form:"roles"`
  }

  SessionV0 struct {
    SessionID string `json:"sessionID" form:"-"`
    Values string `json:"values" form:"-"`
  }
)
