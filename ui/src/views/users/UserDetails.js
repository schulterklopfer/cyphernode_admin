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
            { appList.map( (appData) => (
              <CCard key={appData.id}>
                <CCardHeader>{appData.name}</CCardHeader>
                <CCardBody>
                  { userData.roles.filter( (r) => r.appId === appData.id ).map( (role) => (<CBadge key={role.id} className="font-lg p-2 cna-role text-white" shape="pill">{role.localName}</CBadge>) ) }
                </CCardBody>
              </CCard>
            ) ) }
          </CCardBody>
        </CCollapse>
      </CCard>
    </CCol>
  )
}

export default UserDetails
