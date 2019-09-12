package handlers

var ALLOWED_APP_PROPERTIES = [3]string{ "id", "name", "hash" }
var ALLOWED_USER_PROPERTIES = [4]string{ "id", "name","login","email_address" }

type PagingParams struct {
  Page uint `form:"_page"`
  Limit int `form:"_limit"`
  Sort string `form:"_sort"`
  Order string `form:"_order"`
}

type PagedResult struct {
  Page uint `json:"page"`
  Limit int `json:"limit"`
  Sort string `json:"sort"`
  Order string `json:"order"`
  Total uint `json:"total"`
  Data interface{} `json:"data"`
  Errors []string `json:"errors"`
}
