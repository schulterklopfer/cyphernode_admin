import React from 'react';
import { Redirect } from "react-router-dom";

class AccessMiddleware {

  constructor( ) {
    this.checks = {
      alreadyLoggedIn: this.alreadyLoggedIn.bind(this),
      privateRoute: this.privateRoute.bind(this),
    }
    this.accessPolicies = [];
    this.privateRoutePatterns = [/.+/];
  }

  setAccessPolicies( accessPolicies ) {
    this.accessPolicies = accessPolicies;
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

    if ( !session || !session.token ) {
      // no valid session
      return this._getRouteReturn(true, component);
    }
    return this._getRouteReturn(false, <Redirect to="/"/> );
  }

  privateRoute (component, pathname = '/', session) {

    if ( !session || !session.token ) {
      // no valid session
      return this._getRouteReturn(false, <Redirect to="/login" /> );
    }

    let accessDenied = true;
    for ( const pattern of this.privateRoutePatterns ) {
      if ( pathname.match(pattern) && session && session.token ) {
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
