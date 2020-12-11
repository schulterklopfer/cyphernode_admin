import React from 'react'
import {
  CCard,
  CCardBody,
  CCol,
  CProgress,
  CRow,
} from '@coreui/react'

const Dashboard = () => {
  return (
    <>
      <CCard>
        <CCardBody>
          <CRow className="text-center">
            <CCol md sm="12" className="mb-sm-2 mb-0">
              <div className="text-muted">IBD Progress</div>
              <strong>423453 blocks (60%)</strong>
              <CProgress
                className="progress-xs mt-2"
                precision={1}
                color="success"
                value={60}
              />
            </CCol>
          </CRow>
        </CCardBody>
      </CCard>
    </>
  )
}

export default Dashboard
