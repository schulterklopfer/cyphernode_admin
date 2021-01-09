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
