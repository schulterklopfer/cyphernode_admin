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

import React from 'react'
import {
  CBadge,
  CCard,
  CCardBody,
  CCardHeader,
  CCol,
  CCollapse,
  CLink,
} from '@coreui/react'
import {CIcon} from "@coreui/icons-react";

const UserDetails = (data) => {

  const appList = data.appList;
  const userData = data.userData;

  const [collapsed, setCollapsed] = React.useState(false);

  return (
    <CCol xl={6}>
      <CCard>
        <CCardHeader>
          <CIcon name="cil-user" /> { userData.login } ({userData.email_address})
          <div className="card-header-actions">
            {
              userData.id !== 1 && (
              <CLink className="card-header-action" onClick={() => data.onEditClick(userData) }>
                <CIcon name="cil-pencil" />
              </CLink>)
            }

            <CLink className="card-header-action" onClick={() => setCollapsed(!collapsed)}>
              <CIcon name={collapsed ? 'cil-chevron-top':'cil-chevron-bottom'} />
            </CLink>
          </div>
        </CCardHeader>
        <CCollapse show={collapsed}>
          <CCardBody>
            <h5>Roles</h5>
            { !!userData.description  && <p  class="lead"><em>{ userData.description}</em></p> }
            { appList.map( appData => {
              const filteredRoles = userData.roles.filter( (r) => r.appId === appData.id );
              if ( filteredRoles.length > 0 ) {
                return (
                  <CCard key={appData.id}>
                    <CCardHeader>{appData.name}</CCardHeader>
                    <CCardBody>
                      { filteredRoles.map( (role) => (<CBadge key={role.id} className="font-lg p-2 cna-role text-white mr-1" shape="pill">{role.name}</CBadge>) ) }
                    </CCardBody>
                  </CCard>
                );
              }
              return (<></>)
            } ) }
          </CCardBody>
        </CCollapse>
      </CCard>
    </CCol>
  )
}

export default UserDetails
