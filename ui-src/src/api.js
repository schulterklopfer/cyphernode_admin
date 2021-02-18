
/*
 * MIT License
 *
 * Copyright (c) 2021 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

const debug = 0;

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

  let body = undefined;

  if ( options.body ) {
    if ( typeof options.body === 'object' ) {
      body = JSON.stringify(options.body);
      headers['Content-Type'] = 'application/json';
    } else {
      body = options.body.toString();
    }
  }

  if ( options.session && options.session.token ) {
    headers['Authorization'] = 'Bearer ' + options.session.token;
  }

  return fetch(url, {method, headers, body, redirect: 'error'});
}
