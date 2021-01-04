import React, { useState, useEffect } from 'react'
import {
  CImg,
  CCard,
  CCardHeader,
  CCardBody,
  CCol,
  CProgress,
  CRow, CBadge, CWidgetIcon, CCardFooter, CLink,
} from '@coreui/react'
import requests from "../../requests";
import CIcon from "@coreui/icons-react";

const Dashboard = () => {
  const [status, setStatus] = useState({})
  const [appList, setAppList] = useState({ data: [] })

  useEffect( () => {
    async function fetchStatus() {
      const response = await requests.getStatus();
      if ( response.status === 200 ) {
        // everything is ok
        console.log( "got status" );
        setStatus(response.body);
      }
    }

    fetchStatus();
    const interval = setInterval(fetchStatus, 1000);
    return () => {
      clearInterval(interval);
    }
  }, [] )

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
      <CCard>
        <CCardHeader>
          Bitcoin blockchain <CBadge shape="pill" color="primary">{status.blockchainInfo?.chain}</CBadge> { status && status.blockchainInfo && status.blockchainInfo.initialblockdownload? (
          <div className="card-header-actions">
            <CBadge color="warning" className="float-right">Initial block download</CBadge>
          </div>
        ): (
          <div className="card-header-actions">
            <CBadge color="success" className="float-right">In sync</CBadge>
          </div>
        ) }
        </CCardHeader>
        <CCardBody>
          <CRow>
            <CCol md sm="12" className="mb-sm-2 mb-0">
              <table className="table table-borderless">
              {
                [
                  { field: "blocks", title:"Block count", transform: input => input },
                  { field: "headers", title:"Header count", transform: input => input },
                  { field: "bestblockhash", title:"Best block hash", transform: input => input },
                  { field: "mediantime", title:"Median time", transform: input => {return new Date(input*1000).toLocaleString()} },
                  { field: "difficulty", title:"Difficulty", transform: input => input },
                  { field: "pruned", title:"Pruned", transform: input => { return input?"Yes":"No" } },
                  { field: "verificationprogress", title:"Verfification progress", transform: input => { return (input*100).toFixed(2)+"%" } },
                  { field: "chainwork", title:"Chainwork", transform: input => input }
                ].map( item => (
                  <tr><td className="p-0 m-0"><strong className="nowrap">{ item.title }:</strong></td><td className="p-0 m-0">{ status.blockchainInfo && <span>{ item.transform(status.blockchainInfo[item.field]) } </span> }</td></tr>
                ))
              }
              </table>
            </CCol>
          </CRow>
        </CCardBody>
      </CCard>
      {
        status.blockchainInfo && status.blockchainInfo.initialblockdownload && (
          <CCard>
            <CCardBody>
              <CRow className="text-center">
                <CCol md sm="12" className="mb-sm-2 mb-0">
                  <div className="text-muted">IBD Progress</div>
                  <strong>{status.blockchainInfo.blocks} blocks ({(status.blockchainInfo.verificationprogress*100).toFixed(2)+"%"})</strong>
                  <CProgress
                    className="progress-xs mt-2"
                    precision={1}
                    color="success"
                    value={status.blockchainInfo.verificationprogress*100}
                  />
                </CCol>
              </CRow>
            </CCardBody>
          </CCard>
        )
      }
      <CCard>
        <CCardHeader>
          Intalled Apps
        </CCardHeader>
        <CCardBody>
          <div className="d-flex flex-row flex-wrap justify-content-center">
            {
              appList.data.map((appData) => (
                <CWidgetIcon
                  className="m-2"
                  style={{width:"250px"}}
                  text={appData.version}
                  header={appData.name}
                  iconPadding={false}
                  footerSlot={
                    <CCardFooter className="card-footer px-3 py-2">
                      <CLink
                        className="font-weight-bold font-xs btn-block text-muted"
                        href={"/"+appData.mountPoint}
                        rel="noopener norefferer"
                        target="_blank"
                      >
                        Go to app
                        <CIcon name="cil-arrow-right" className="float-right" width="16"/>
                      </CLink>
                    </CCardFooter>
                  }
                >
                  <CImg className="p-2" style={{ "background-color": appData.meta?.color, "border-radius":"10px"}} width="60" height="60" src={appData.meta?.icon}/>
                </CWidgetIcon>
              ))
            }
          </div>
        </CCardBody>
      </CCard>
    </>
  )
}

export default Dashboard
