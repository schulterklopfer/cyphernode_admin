import React from 'react';
import * as Cookies from "js-cookie";

class Session {

  constructor() {
    this.token = this._getCookie();
    this.payload = this._parseJWT(this.token);
  }

  save() {
    this._setCookie(this.token);
  }

  _parseJWT( token ) {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload);
  }

  _setCookie(session) {
    Cookies.remove("session");
    // sameSite: "Strict", HttpOnly: true
    Cookies.set("session", session, { expires: 14});
  }

  _getCookie() {
    return Cookies.get("session");
  }

  setToken( token ) {
    this.token = token;
  }

  getToken() {
    return this.token;
  }
}

export const getSession = () => { return new Session() };
const SessionContext = React.createContext( getSession() );
export default SessionContext;
