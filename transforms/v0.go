package transforms

type (
  AppV0 struct {
    ID uint   `json:"id" form:"id"`
    ClientID string `json:"clientID" form:"clientid_like"`
    Name string `json:"name" form:"name_like"`
    Description string `json:"description" form:"description"`
    CallbackURL string `json:"callbackURL" form:"callbackURL"`
    PostLogoutCallbackURL string `json:"postLogoutCallbackURL" form:"callbackURL"`
    AvailableRoles []*RoleV0 `json:"availableRoles" form:"availableRoles"`
  }

  RoleV0 struct {
    ID uint   `json:"id" form:"id"`
    Name string `json:"name" form:"name"`
    Description string `json:"description" form:"description"`
    AutoAssign bool `json:"autoAssign" form:"autoAssign"`
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
