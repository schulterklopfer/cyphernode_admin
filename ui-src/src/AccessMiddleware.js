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
import { Redirect } from "react-router-dom";
import { getAdminApp } from "./redux/selectors";

class AccessMiddleware {

  constructor( ) {
    this.checks = {
      alreadyLoggedIn: this.alreadyLoggedIn.bind(this),
      privateRoute: this.privateRoute.bind(this),
    }
    this.privateRoutePatterns = [/.+/];
  }

  getAccessPolicies() {
    const adminApp = getAdminApp();
    if ( adminApp ) {
      return adminApp.accessPolicies;
    }
    return [];
  }

  routeDisplay (middlewares= [], routeToVisit, directedFrom = '', session = {}) {
    let ret = null;
    for (let i = 0; i < middlewares.length; i++) {
      if ( !this.checks[middlewares[i]] ) {
        console.log( "Unknow access check "+middlewares[i]);
        continue;
      }
      ret = this.checks[middlewares[i]](routeToVisit, directedFrom, session);
      if (!ret.evalNext) {
        break;
      }
    }
    if ( ret ) {
      return ret.routeObject;
    }
    return this._getRouteReturn(false, <Redirect to="/404" /> );
  }

  _getRouteReturn (evalNext, routeObject) {
    return {evalNext, routeObject};
  }

  alreadyLoggedIn (component, _, session) {

    if ( !session || !session.token || !session.jwt || !session.jwt.valid) {
      // no valid session
      return this._getRouteReturn(true, component);
    }
    return this._getRouteReturn(false, <Redirect to="/"/> );
  }

  privateRoute (component, pathname = '/', session) {

    if ( !session || !session.token || !session.jwt.valid ) {
      // no valid session
      return this._getRouteReturn(false, <Redirect to="/login" /> );
    }

    let accessDenied = true;
    for ( const pattern of this.privateRoutePatterns ) {
      if ( pathname.match(pattern) && session && session.token && session.jwt.valid) {
        accessDenied = false;
      }
      if ( !accessDenied ) {
        break;
      }
    }

    if( accessDenied ) {
      return this._getRouteReturn(false, <Redirect to="/404" /> );
    }
    return this._getRouteReturn(true, component);
  }
}

export default AccessMiddleware;
