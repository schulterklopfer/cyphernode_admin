package transforms

import "github.com/schulterklopfer/cyphernode_admin/models"

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
