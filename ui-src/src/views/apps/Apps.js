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
