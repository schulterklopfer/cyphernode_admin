//const apiBaseUri = 'http://192.168.178.90/admin/api/v0';
//const apiBaseUri = '192.168.178.90:3030/api/v0';
const webScheme = 'http';
const webSocketScheme = 'ws';
const apiBaseUri = '/admin/api/v0';

const requests = {
  login: async (username, password ) => {
    return await requests.request(
      'POST',
      webScheme+"://"+apiBaseUri+'/login',
      { body: { username, password } } );
  },

  getStatus: async ( session ) => {
    return await requests.request(
      'GET',
      webScheme+"://"+apiBaseUri+'/status',
      { session }
    );
  },

  getDockerContainerByName: async (name, session ) => {
    return await requests.request(
      'GET',
      webScheme+"://"+apiBaseUri+'/docker/name/'+name,
      { session }
    );
  },

  getDockerContainerByImageHash: async ( hash, session ) => {
    return await requests.request(
      'GET',
      webScheme+"://"+apiBaseUri+'/docker/image/'+hash,
      { session }
    );
  },

  getDockerLogsWebsocketClient: ( containerId, session ) => {
    return new WebSocket(webSocketScheme+"://"+apiBaseUri+'/docker/logs/'+containerId );
  },

  getUsers: async ( session ) => {
    return await requests.request(
      'GET',
      webScheme+"://"+apiBaseUri+'/users',
      { session }
    );
  },

  createUser: async ( user, session ) => {
    return await requests.request(
      'POST',
      webScheme+"://"+apiBaseUri+'/users',
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
      webScheme+"://"+apiBaseUri+'/users/'+id,
      { session, body: user }
    );
  },

  deleteUser: async ( id, session ) => {
    return await requests.request(
      'DELETE',
      webScheme+"://"+apiBaseUri+'/users/'+id,
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
      webScheme+"://"+apiBaseUri+'/apps/'+id,
      { session }
    );
  },

  getApps: async ( session ) => {
    return await requests.request(
      'GET',
      webScheme+"://"+apiBaseUri+'/apps',
      { session }
    );
  },

  createApp: async ( app, session ) => {
    return await requests.request(
      'POST',
      webScheme+"://"+apiBaseUri+'/apps',
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
      webScheme+"://"+apiBaseUri+'/apps/'+id,
      { session, body: app }
    );
  },

  deleteApp: async (id, session ) => {
    return await requests.request(
      'DELETE',
      webScheme+"://"+apiBaseUri+'/apps/'+id,
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

    const response = await fetch(url, {method, headers, body});

    const r = {
      status: response.status,
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
