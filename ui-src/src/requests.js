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

const wsBaseUri =   (protocol==='https:'?'wss:':'ws:')+'//' + host + mountPoint + '/api/v0';
const apiBaseUri =  protocol + '//' + host + mountPoint + '/api/v0';


const requests = {
  login: async (username, password ) => {
    return await requests.request(
      'POST',
      apiBaseUri+'/login',
      { body: { username, password } } );
  },

  getStatus: async ( session ) => {
    return await requests.request(
      'GET',
      apiBaseUri+'/status',
      { session }
    );
  },

  getMe: async ( session ) => {
    return await requests.request(
      'GET',
      apiBaseUri+'/users/me',
      { session }
    );
  },

  getDockerContainerByName: async (name, session ) => {
    return await requests.request(
      'GET',
      apiBaseUri+'/docker/name/'+name,
      { session }
    );
  },

  getDockerContainerByImageHash: async ( hash, session ) => {
    return await requests.request(
      'GET',
      apiBaseUri+'/docker/image/'+hash,
      { session }
    );
  },

  getDockerLogsWebsocketClient: ( containerId, session ) => {
    return new WebSocket(wsBaseUri+'/docker/logs/'+containerId );
  },

  getUsers: async ( session ) => {
    return await requests.request(
      'GET',
      apiBaseUri+'/users',
      { session }
    );
  },

  createUser: async ( user, session ) => {
    return await requests.request(
      'POST',
      apiBaseUri+'/users',
      { session, body: user }
    );
  },

  patchUser: async ( id, user, session ) => {

    if( !id ) {
      return {
        status: -1
      }
    }

    return await requests.request(
      'PATCH',
      apiBaseUri+'/users/'+id,
      { session, body: user }
    );
  },

  deleteUser: async ( id, session ) => {
    return await requests.request(
      'DELETE',
      apiBaseUri+'/users/'+id,
      { session }
    );
  },

  getApp: async ( id, session ) => {
    if( !id ) {
      return {
        status: 404
      }
    }

    return await requests.request(
      'GET',
      apiBaseUri+'/apps/'+id,
      { session }
    );
  },

  getApps: async ( session ) => {
    return await requests.request(
      'GET',
      apiBaseUri+'/apps',
      { session }
    );
  },

  createApp: async ( app, session ) => {
    return await requests.request(
      'POST',
      apiBaseUri+'/apps',
      { session, body: app }
    );
  },

  patchApp: async ( id, app, session ) => {

    if( !id ) {
      return {
        status: 404
      }
    }

    return await requests.request(
      'PATCH',
      apiBaseUri+'/apps/'+id,
      { session, body: app }
    );
  },

  deleteApp: async (id, session ) => {
    return await requests.request(
      'DELETE',
      apiBaseUri+'/apps/'+id,
      { session }
    );
  },

  request: async ( method, url, options ) => {

    options = options || {};

    const headers = Object.assign({}, options.headers );

    if ( options.session && options.session.token ) {
      headers['Authorization'] = 'Bearer ' + options.session.token;
    }

    let body = undefined;

    if ( options.body ) {
      if ( typeof options.body === 'object' ) {
        body = JSON.stringify(options.body);
        headers['Content-Type'] = 'application/json';
      } else {
        body = options.body.toString();
      }
    }

    let response;

    try {
      // will throw when redirected
      response = await fetch(url, {method, headers, body, redirect: 'error'});
    } catch( err ) {
      console.log( err );
      return;
    }

    const r = {
      status: response.status,
      headers: response.headers
    }

    const ct = response.headers.get('Content-Type');

    if ( ct && ct.startsWith('application/json') ) {
      r.body = await response.json();
    } else {
      r.body = await response.text();
    }
    return r;
  }
};

export default requests;
