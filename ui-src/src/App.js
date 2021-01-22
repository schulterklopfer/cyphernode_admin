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

import React, {Component} from 'react';
import {HashRouter, Route, Switch} from 'react-router-dom';
import './scss/style.scss';
import SessionContext, {getSession} from "./sessionContext";
import AccessMiddleware from "./AccessMiddleware";
import { ErrorCodes } from "./errors";
import requests from "./requests";


const loading = (
  <div className="pt-3 text-center">
    <div className="sk-spinner sk-spinner-pulse"></div>
  </div>
)
// https://www.jmfurlott.com/handling-user-session-react-context/
// Containers
const TheLayout = React.lazy(() => import('./containers/TheLayout'));

// Pages
const Login = React.lazy(() => import('./views/pages/login/Login'));
const Logout = React.lazy(() => import('./views/pages/logout/Logout'));
const Register = React.lazy(() => import('./views/pages/register/Register'));
const Page404 = React.lazy(() => import('./views/pages/page404/Page404'));
const Page500 = React.lazy(() => import('./views/pages/page500/Page500'));

class App extends Component {

  constructor( props ) {
    super( props );
    this.login = this.login.bind(this);
    this.accessMiddleware = new AccessMiddleware();
    this.state = {
      login: this.login
    };
  }

  async componentDidMount() {
    const session = await getSession();
    this.setState( { session } );
    let app = {};
    try {
      const appResponse = await requests.getApp(1, this.state.session );
      if ( appResponse.status === 200 ) {
        app = appResponse.body;
        this.accessMiddleware.setAccessPolicies(app.accessPolicies);
        this.setState({app:app});
      }
    } catch (e) {
      console.log( e );
    }
  }

  async login(username, password ) {
    try {
      const response = await requests.login( username, password );
      console.log( response );
      if ( response.status === 200 ) {
        // everything is ok
        await this.state.session.setToken( response.body.token );
        this.state.session.save();
      } else if( response.status === 401 ) {
        return ErrorCodes.WRONG_CREDENTIALS;
      } else {
        return ErrorCodes.UNKNOWN_ERROR;
      }
    } catch (e) {
      console.log(e);
      return ErrorCodes.API_UNAVAILABLE;
    }
    return 0;
  }

  render() {
    return (
      <HashRouter>
        <React.Suspense fallback={loading}>
          <SessionContext.Provider value={this.state}>
            <Switch>
              <Route exact path="/login" name="Login" render={
                props => this.accessMiddleware.routeDisplay(
                  ['alreadyLoggedIn'],
                  <Login {...props}/>,
                  props.location.pathname,
                  this.state.session )
              }/>

              <Route exact path="/logout" name="Logout" render={props => <Logout {...props}/>} />
              <Route exact path="/register" name="Register Page" render={props => <Register {...props}/>} />
              <Route exact path="/404" name="Page 404" render={props => <Page404 {...props}/>} />
              <Route exact path="/500" name="Page 500" render={props => <Page500 {...props}/>} />

              <Route path="/" name="Home" render={
                props => this.accessMiddleware.routeDisplay(
                  ['privateRoute'],
                  <TheLayout {...props}/>,
                  props.location.pathname,
                  this.state.session )
              }/>
            </Switch>
          </SessionContext.Provider>
        </React.Suspense>
      </HashRouter>
    );
  }
}

export default App;
