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

import React, {useContext} from 'react'
import {
  CCard,
  CCardBody,
  CCardHeader,
  CCol,
} from '@coreui/react'
import SessionContext from "../../sessionContext";

const UserProfile = () => {

  const context = useContext( SessionContext );

  return (
    <CCol xl={6}>
      <CCard>
        <CCardHeader className="h5">
          Your profile
        </CCardHeader>
        <CCardBody>
          <pre>
          { JSON.stringify( context.session.user, null, 2 ) }
          </pre>
          <table className="table table-borderless p-0 m-0">
            <tbody>
            <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Login:</td><td className="p-0 pr-1 pl-1 m-0">{ context.session.user.login }</td></tr>
            <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">name:</td><td className="p-0 pr-1 pl-1 m-0">{ context.session.user.name }</td></tr>
            <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Email address:</td><td className="p-0 pr-1 pl-1 m-0">{ context.session.user.email_address }</td></tr>
            <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Roles:</td><td className="p-0 pr-1 pl-1 m-0">
              { context.session.user.roles.map( role => (
                <table key={role.id} className="table table-borderless p-0 m-0">
                  <tbody>
                  <tr><td className="p-0 pr-1 pl-1 m-0">name:</td><td className="p-0 pr-1 pl-1 m-0"></td></tr>
                  <tr><td className="p-0 pr-1 pl-1 m-0">description:</td><td className="p-0 pr-1 pl-1 m-0"></td></tr>
                  <tr><td className="p-0 pr-1 pl-1 m-0">App:</td><td className="p-0 pr-1 pl-1 m-0"></td></tr>
                  </tbody>
                </table>
              ) ) }
            </td></tr>
            </tbody>
          </table>
        </CCardBody>
      </CCard>
    </CCol>
  )
}

export default UserProfile
