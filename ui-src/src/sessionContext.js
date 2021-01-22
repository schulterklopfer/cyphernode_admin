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

import React from 'react';
import * as Cookies from "js-cookie";
import requests from "./requests";

class Session {

  constructor() {
    const token = this._getCookie();
    this.setToken( token );
  }

  async appendLocalData() {
    // load user from database and put it in session.
    if ( this.payload.id ) {
      // we have an id
      const response = await requests.getMe( this );
      if ( response && response.status === 200 ) {
        // all good
        this.user = response.body;
        return true;
      }
    }
    this.end();
    return false;
  }

  save() {
    this._setCookie(this.token);
  }

  end() {
    this.setToken( undefined );
    this._removeCookie();
  }

  _parseJWT( token ) {
    if ( !token || token === ''  ) {
      return {};
    }
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    const jwt = JSON.parse(jsonPayload);
    jwt.valid = this._JWTIsValid(jwt);
    return jwt;
  }

  _JWTIsValid( jwt ) {
    const now = ((new Date()).getTime()/1000)<<0;

    if( !this._verifyExp( jwt.exp, now, true ) ) {
      return false;
    }

    if( !this._verifyIat( jwt.iat, now, false ) ) {
      return false;
    }

    return true;
  }

  _verifyExp(exp, now, req) {
    if (!exp) {
      return !req;
    }
    return now <= exp;
  }

  _verifyIat(iat, now, req ) {
    if (!iat) {
      return !req;
    }
    return now >= iat;
  }

  _setCookie(session) {
    Cookies.remove("session");
    // sameSite: "Strict", HttpOnly: true
    Cookies.set("session", session, { expires: 14});
  }

  _getCookie() {
    return Cookies.get("session");
  }

  _removeCookie() {
    Cookies.remove("session");
  }

  async setToken( token ) {
    this.token = token;
    if( token ) {
      this.jwt = this._parseJWT( token );
      this.payload = this._parseJWT(this.token);
      await this.appendLocalData();
    } else {
      delete this.jwt
      delete this.payload;
    }

  }

  getToken() {
    return this.token;
  }
}

export const getSession = async () => {
  const session = new Session();
  return session
};

const SessionContext = React.createContext( getSession() );
export default SessionContext;
