
const debug = 1;

let host = window.location.host;
let protocol = window.location.protocol;
let mountPoint = "/admin";

if ( debug === 1) {
  host = "localhost:3030";
  protocol = "http:";
  mountPoint = "";
} else if ( debug === 2 ) {
  host = "localhost";
  protocol = "https:";
}

//const wsBaseUri =   (protocol==='https:'?'wss:':'ws:')+'//' + host + mountPoint + '/api/v0';
const apiBaseUri =  protocol + '//' + host + mountPoint + '/api/v0';
const wsBaseUri =   (protocol==='https:'?'wss:':'ws:')+'//' + host + mountPoint + '/api/v0';

export function __fetchApps( session ) {
  return apiFetch(
    'GET',
    apiBaseUri+'/apps',
    { session }
  );
}

export function __fetchStatus( session ) {
  return apiFetch(
    'GET',
    apiBaseUri + '/status',
    {session}
  );
}

export function __fetchUsers( session ) {
  return apiFetch(
    'GET',
    apiBaseUri+'/users',
    { session }
  );
}

export function __createUser(userData, session ) {
  return apiFetch(
    'POST',
    apiBaseUri+'/users',
    { session, body: userData }
  );
}

export function __patchUser( userId, userData, session ) {
  return apiFetch(
    'PATCH',
    apiBaseUri+'/users/'+userId,
    { session, body: userData }
  );
}

export function __deleteUser( userId, session ) {
  return apiFetch(
    'DELETE',
    apiBaseUri+'/users/'+userId,
    { session }
  );
}

export function __fetchContainerInfoByName( name, session ) {
  return apiFetch(
    'GET',
    apiBaseUri+'/docker/name/'+name,
    { session }
  );
}

export function __fetchContainerInfoByHash( hash, session ) {
  return apiFetch(
    'GET',
    apiBaseUri+'/docker/image/'+hash,
    { session }
  );
}

export function login(username, password ) {
  return apiFetch(
    'POST',
    apiBaseUri+'/login',
    { body: { username, password } }
  );
}

export function resolveFile( file ) {
  return apiBaseUri+'/files/'+file;
}

export function fetchDockerLogsWebsocketClient( containerId, session ) {
  return new WebSocket(wsBaseUri+'/docker/logs/'+containerId );
}

async function apiFetch( method, url, options ) {

  options = options || {};

  const headers = Object.assign({}, options.headers );

  if ( options.session && options.session.token ) {
    headers['Authorization'] = 'Bearer ' + options.session.token;
  }

  return fetch(url, {method, headers, redirect: 'error'});
}
