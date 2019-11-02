package globals

const SESSION_COOKIE_NAME string = "_cna_session"

/** env keys **/
const HYDRA_ADMIN_URL_ENV_KEY string = "CNA_HYDRA_ADMIN_URL"
const OIDC_SESSION_COOKIE_SECRET_ENV_KEY string = "OIDC_SESSION_COOKIE_SECRET"
const HYDRA_DISABLE_SYNC_ENV_KEY string = "CNA_DISABLE_HYDRA_SYNC"
const OIDC_DISCOVERY_URL_ENV_KEY string = "OIDC_DISCOVERY_URL"
const BASE_URL_ENV_KEY string = "BASE_URL"
const OIDC_SSO_COOKIE_DOMAIN_ENV_KEY string = "OIDC_SSO_COOKIE_DOMAIN"
const CNA_ADMIN_LOGIN_ENV_KEY string = "CNA_ADMIN_LOGIN"
const CNA_ADMIN_PASSWORD_ENV_KEY string = "CNA_ADMIN_PASSWORD"
const CNA_ADMIN_NAME_ENV_KEY string = "CNA_ADMIN_NAME"
const CNA_ADMIN_EMAIL_ADDRESS_ENV_KEY string = "CNA_ADMIN_EMAIL_ADDRESS"
const CNA_ADMIN_DATABASE_FILE_ENV_KEY string = "CNA_ADMIN_DATABASE_FILE"
const CNA_ADMIN_APP_WHITELIST_FILE_ENV_KEY = "CNA_ADMIN_APP_WHITELIST_FILE"

/** defaults **/
const DEFAULTS_BASE_URL string = "http://127.0.0.1:3030"
const DEFAULTS_OIDC_DISCOVERY_URL string = "http://127.0.0.1:9000/.well-known/openid-configuration"
const DEFAULTS_OIDC_SESSION_COOKIE_SECRET string = "secret"
const DEFAULTS_OIDC_SSO_COOKIE_DOMAIN string = "127.0.0.1"
const DEFAULTS_CNA_ADMIN_LOGIN string = "admin"
const DEFAULTS_CNA_ADMIN_PASSWORD string = "admin"
const DEFAULTS_CNA_ADMIN_NAME string = "admin"
const DEFAULTS_CNA_ADMIN_EMAIL_ADDRESS string = "admin@admin.com"
const DEFAULTS_CNA_ADMIN_DATABASE_FILE string = "/data/db.sqlite3"
const DEFAULTS_CNA_ADMIN_APP_WHITELIST_FILE string = "./apps.txt"

/** general hydra stuff **/
const HYDRA_SCOPE_OFFLINE string = "offline"
const HYDRA_SCOPE_OPEN_ID string = "openid"

/** router group names **/
const ROUTER_GROUPS_PUBLIC string = "public"
const ROUTER_GROUPS_PRIVATE string = "private"
const ROUTER_GROUPS_HYDRA string = "hydra"
const ROUTER_GROUPS_SESSIONS string = "sessions"
const ROUTER_GROUPS_USERS string = "users"
const ROUTER_GROUPS_APPS string = "apps"

/** router groups endpoint bases **/
const ROUTER_GROUPS_BASE_ENDPOINT_PUBLIC string = ""
const ROUTER_GROUPS_BASE_ENDPOINT_PRIVATE string = "/_"
const ROUTER_GROUPS_BASE_ENDPOINT_HYDRA string = "/hydra"
const ROUTER_GROUPS_BASE_ENDPOINT_SESSIONS string = "/sessions"
const ROUTER_GROUPS_BASE_ENDPOINT_USERS string = "/api/v0/users"
const ROUTER_GROUPS_BASE_ENDPOINT_APPS string = "/api/v0/apps"

/** urls and endpoints **/
const PUBLIC_ENDPOINTS_LOGIN string = "/login"
const PUBLIC_ENDPOINTS_CALLBACK string = "/callback"
const PUBLIC_ENDPOINTS_BYEBYE string = "/byebye"
const PRIVATE_ENDPOINTS_HOME string = "/home"

//const URLS_OIDC_DISCOVERY string = "http://127.0.0.1:9000/.well-known/openid-configuration"
const URLS_CALLBACK string = ROUTER_GROUPS_BASE_ENDPOINT_PUBLIC+PUBLIC_ENDPOINTS_CALLBACK

const URLS_BYEBYE string =  ROUTER_GROUPS_BASE_ENDPOINT_PUBLIC+PUBLIC_ENDPOINTS_BYEBYE


/** sql statements **/
const SQL_STATEMENTS__ROLES_BY_USER_ID_AND_APP_ID string =
    "SELECT " +
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
// add offline scope, so we get refresh tokens
var HYDRA_AUTO_SCOPES = [...]string{ HYDRA_SCOPE_OFFLINE, HYDRA_SCOPE_OPEN_ID }

var ENDPOINTS_PUBLIC_PATTERNS = [...]string{ ".*/+favicon.ico$" }

var ROUTER_GROUPS = [...]string{
  /* public */
  ROUTER_GROUPS_PUBLIC,
  ROUTER_GROUPS_HYDRA,
  ROUTER_GROUPS_SESSIONS,
  /* protected */
  ROUTER_GROUPS_PRIVATE,
  ROUTER_GROUPS_APPS,
  ROUTER_GROUPS_USERS,
}
var ROUTER_GROUPS_BASE_ENDPOINTS = [...]string{
  /* public */
  ROUTER_GROUPS_BASE_ENDPOINT_PUBLIC,
  ROUTER_GROUPS_BASE_ENDPOINT_HYDRA,
  ROUTER_GROUPS_BASE_ENDPOINT_SESSIONS,
  /* protected */
  ROUTER_GROUPS_BASE_ENDPOINT_PRIVATE,
  ROUTER_GROUPS_BASE_ENDPOINT_APPS,
  ROUTER_GROUPS_BASE_ENDPOINT_USERS,
}
var PROTECTED_ROUTER_GROUPS_INDICES = [...]int{ 3,4,5 }
