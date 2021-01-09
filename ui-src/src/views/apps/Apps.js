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

import React, { useState, useEffect } from 'react'
import {
  CCard,
  CCardBody,
  CRow, CCol
} from '@coreui/react'
import AppDetails from "./AppDetails";
import requests from "../../requests";

const Apps = () => {
  const [appList, setAppList] = useState({ data: [] })

  useEffect( () => {
    let ignore = false;
    async function fetchAppList() {
      const response = await requests.getApps();
      if ( response.status === 200 ) {
        // everything is ok
        if ( !ignore ) {
          setAppList(response.body);
        }
      }
    }

    fetchAppList();
    return () => { ignore = true }
  }, [] )

  return (
    <>
    <CRow><CCol>
      <CCard>
        <CCardBody>
          <div className="mb-3">
          The following cypherapps have been installed using the <code>cam</code> tool. They
          cannot (yet) be edited using the cyphernode admin. It's just a quick overview
          of all the roles a cypherapp provides and how access to different resources
          is handled.
          </div>
          <div className="mb-3">
            <h5>Available Roles</h5>
            The available roles describe the roles present in the respective cypherapp. Please note
            that all roles are local and for example the <code>admin</code> role in one cypherapp is not
            the same as the <code>admin</code> role in another cypherapp. If a role has
            the <code>autoAssign</code> flag set, it means that every user will automatically receive
            that role. This is useful, when you want to install a new cypherapp and give the existing
            users access to the regular user space of that cypherapp.
          </div>
          <div className="mb-3">
            <h5>Access Policies</h5>
            Access policies allow or deny access of certain roles within the cypherapp
            to a certain resource relative to the mount point of the cypherapp. All
            resources are described using regular expressions. The actions are equivalent
            to the available HTTP methods.
          </div>
        </CCardBody>
      </CCard>
    </CCol></CRow>
    <CRow> { appList.data.map((appData) => ( <AppDetails key={appData.id} {...appData} /> ) ) } </CRow>
    </>
  )
}

export default Apps
