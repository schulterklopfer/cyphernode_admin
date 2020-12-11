package globals

const SESSION_COOKIE_NAME string = "cna_session"

/** env keys **/
const BASE_URL_EXTERNAL_ENV_KEY string = "BASE_URL_EXTERNAL"
const BASE_URL_INTERNAL_ENV_KEY string = "BASE_URL_INTERNAL"
const OIDC_SSO_COOKIE_DOMAIN_ENV_KEY string = "OIDC_SSO_COOKIE_DOMAIN"
const CNA_ADMIN_LOGIN_ENV_KEY string = "CNA_ADMIN_LOGIN"
const CNA_ADMIN_PASSWORD_ENV_KEY string = "CNA_ADMIN_PASSWORD"
const CNA_ADMIN_NAME_ENV_KEY string = "CNA_ADMIN_NAME"
const CNA_ADMIN_EMAIL_ADDRESS_ENV_KEY string = "CNA_ADMIN_EMAIL_ADDRESS"
const CNA_ADMIN_DATABASE_FILE_ENV_KEY string = "CNA_ADMIN_DATABASE_FILE"
const CYPHERAPPS_INSTALL_DIR_ENV_KEY = "CYPHERAPPS_INSTALL_DIR"
const CNA_COOKIE_SECRET_ENV_KEY = "CNA_COOKIE_SECRET"

/** router groups endpoint bases **/
const BASE_ENDPOINT_PUBLIC string = ""
const BASE_ENDPOINT_USERS string = "/api/v0/users"
const BASE_ENDPOINT_APPS string = "/api/v0/apps"
const BASE_ENDPOINT_SESSIONS string = "/api/v0/sessions"

/** urls and endpoints **/
const FORWARD_AUTH_ENDPOINTS_AUTH = "/forwardauth"

const PUBLIC_ENDPOINTS_LOGIN string = "/api/v0/login"
const PRIVATE_ENDPOINTS_LOGOUT string = "/api/v0/logout"

const CYPHERAPPS_REPO string = "git://github.com/SatoshiPortal/cypherapps.git"

/** sql statements **/
const SQL_STATEMENTS__ROLES_BY_USER_ID_AND_APP_ID string = "SELECT " +
    "role_models.id as id, " +
    "role_models.app_id as app_id, " +
    "role_models.auto_assign as auto_assign, " +
    "role_models.name as name, " +
    "role_models.description as description, " +
    "role_models.created_at as created_at, " +
    "role_models.updated_at as updated_at, " +
    "role_models.deleted_at as deleted_at " +
    "FROM role_models " +
    "JOIN user_roles " +
    "ON role_models.id = user_roles.role_model_id " +
    "WHERE user_roles.user_model_id = ? " +
    "AND role_models.app_id = ?"

/** roles **/
const ROLES_ADMIN_ROLE = "admin"

/** useful vars **/

var ENDPOINTS_PUBLIC_PATTERNS = [...]string{".*/+favicon.ico$"}

var PROTECTED_ROUTER_GROUPS_INDICES = [...]int{2, 3, 4}

/** defaults **/

var DEFAULTS = map[string]string{
  BASE_URL_EXTERNAL_ENV_KEY:       "http://www.cna.localhost:3030",
  BASE_URL_INTERNAL_ENV_KEY:       "http://www.cna.localhost:3031",
  OIDC_SSO_COOKIE_DOMAIN_ENV_KEY:  "www.cna.localhost",
  CNA_COOKIE_SECRET_ENV_KEY:       "thisIsTheDefaultSecret",
  CNA_ADMIN_LOGIN_ENV_KEY:         "admin",
  CNA_ADMIN_PASSWORD_ENV_KEY:      "admin",
  CNA_ADMIN_NAME_ENV_KEY:          "admin",
  CNA_ADMIN_EMAIL_ADDRESS_ENV_KEY: "admin@admin.com",
  CNA_ADMIN_DATABASE_FILE_ENV_KEY: "/data/db.sqlite3",
  CYPHERAPPS_INSTALL_DIR_ENV_KEY:  "/data/apps",
}
