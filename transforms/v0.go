package transforms

type (
  AppV0 struct {
    ID uint   `json:"id" form:"id"`
    Hash string `json:"hash" form:"hash"`
    Name string `json:"name" form:"name"`
    Description string `json:"description" form:"description"`
    AvailableRoles []*RoleV0 `json:"availableRoles" form:"availableRoles"`
  }

  RoleV0 struct {
    ID uint   `json:"id" form:"id"`
    Name string `json:"name" form:"name"`
    Description string `json:"description" form:"description"`
    AutoAssign bool `json:"autoAssign" form:"autoAssign"`
    AppId uint `json:"appId" form:"appId"`
  }

  // strange form names come from ng2_smart_table
  // didn't want to change the defaults
  UserV0 struct {
    ID uint   `json:"id" form:"id_like"`
    Login string `json:"login" form:"login_like"`
    Name string `json:"name" form:"name_like"`
    EmailAddress string `json:"emailAddress" form:"emailAddress_like"`
    Roles []*RoleV0 `json:"roles" form:"roles"`
  }
)
