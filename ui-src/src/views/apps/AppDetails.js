import React from 'react'
import {
  CCard,
  CCardBody,
  CCardHeader, CCol, CDataTable,
} from '@coreui/react'
import {CIcon} from "@coreui/icons-react";

const AppDetails = (appData) => {
  return (
    <CCol xl={6}>
    <CCard>
      <CCardHeader>
        <CIcon name="cil-applications" /> { appData.name } - <code>{appData.hash}</code>
        <div className="float-right">{ '/'+appData.mountPoint }</div>
      </CCardHeader>
      <CCardBody>
        { !!appData.description  && <div className="mb-3 mt-mi lead"><em>{ appData.description}</em></div> }

        <h5>Available Roles</h5>
        <CDataTable
          striped={true}
          items={appData.availableRoles}
          fields={['name','description','autoAssign']}
          size="sm"
        />

        <h5>Access Policies</h5>
        <CDataTable
          striped={true}
          items={appData.accessPolicies}
          fields={['roles','resources','effect','actions']}
          size="sm"
        />
      </CCardBody>
    </CCard>
    </CCol>
  )
}

export default AppDetails
