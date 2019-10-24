package globals

/** env keys **/
const HYDRA_ADMIN_URL_ENV_KEY string = "CNA_HYDRA_ADMIN_URL"
const HYDRA_DISABLE_SYNC_ENV_KEY string = "CNA_DISABLE_HYDRA_SYNC"

/** general hydra stuff **/
const HYDRA_SCOPE_OFFLINE string = "offline"
const HYDRA_SCOPE_OPEN_ID string = "openid"

/** urls and endpoints **/
const URL_HOST string = "http://127.0.0.1:3030"
const ENDPOINTS_LOGIN string = "/public/login"
const ENDPOINTS_CALLBACK string = "/hydra/callback"
const ENDPOINTS_HOME string = "/home"

const URLS_OIDC_DISCOVERY string = "http://127.0.0.1:9000/.well-known/openid-configuration"
const URLS_CALLBACK string = URL_HOST+ENDPOINTS_CALLBACK

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
var ENDPOINTS_PUBLIC_PATTERNS = [...]string{"^/public/", "^/hydra/", ".*/favicon.ico$" }
