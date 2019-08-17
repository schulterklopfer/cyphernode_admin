package transforms

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
