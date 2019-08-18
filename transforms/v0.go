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

  UserV0 struct {
    ID uint   `json:"id" form:"id"`
    Login string `json:"login" form:"login"`
    Name string `json:"name" form:"name"`
    EmailAddress string `json:"emailAddress" form:"emailAddress"`
    Roles []*RoleV0 `json:"roles" form:"roles"`
  }
)
