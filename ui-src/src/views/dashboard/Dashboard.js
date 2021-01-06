import React, { useState, useEffect } from 'react'
import {
  CImg,
  CCard,
  CCardHeader,
  CCardBody,
  CCol,
  CProgress,
  CRow, CBadge, CWidgetIcon, CCardFooter, CLink, CTooltip,
} from '@coreui/react'
import requests from "../../requests";
import CIcon from "@coreui/icons-react";
import crypto from "crypto"

const Dashboard = () => {
  const [status, setStatus] = useState({})
  const [appList, setAppList] = useState({ data: [] })

  useEffect( () => {
    async function fetchStatus() {
      const response = await requests.getStatus();
      if ( response.status === 200 ) {
        // everything is ok
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
          const appList = response.body;
          // do not include admin app here
          appList.data = appList.data.filter( app => app.id!==1 )

          for ( const app of appList.data ) {
            const containerResponse = await requests.getDockerContainerByHash( app.hash );

            if ( containerResponse.status === 200 ) {
              app.container = {
                state: containerResponse.body.State,
                status: containerResponse.body.Status,
                created: containerResponse.body.Created,
              };
            }

          }

          setAppList(appList);
        }
      }
    }

    fetchAppList();
    return () => { ignore = true }
  }, [] )

  return (
    <>
      {
        status.blockchainInfo && status.blockchainInfo.initialblockdownload && (
          <CCard>
            <CCardHeader>
              Initial block download
            </CCardHeader>
            <CCardBody>
              <CRow className="text-center">
                <CCol md sm="12" className="mb-sm-2 mb-0">
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
                  { field: "mediantime", title:"Median time", transform: input => {return new Date(input*1000).toLocaleString()} },
                  { field: "difficulty", title:"Difficulty", transform: input => input },
                  { field: "pruned", title:"Pruned", transform: input => { return input?"Yes":"No" } },
                  { field: "verificationprogress", title:"Verfification progress", transform: input => { return (input*100).toFixed(2)+"%" } },
                ].map( item => (
                  <tr><td className="p-0 m-0"><strong className="nowrap">{ item.title }:</strong></td><td className="p-0 m-0">{ status.blockchainInfo && <span>{ item.transform(status.blockchainInfo[item.field]) } </span> }</td></tr>
                ))
              }
              </table>
            </CCol>
          </CRow>
        </CCardBody>
          <CCardFooter>
            <div className="d-flex flex-row flex-wrap justify-content-center">
              { status.latestBlocks?.map( (block, index) => (
                <div className="d-flex flex-row align-items-center">
                <CTooltip
                  content={ "Hash: "+block.hash }
                  placement="top"
                >
                  <CCard className="m-2 font-xs">
                    <CCardHeader color="primary" className="p-1">
                      <CLink
                        className="font-weight-bold btn-block text-light text-center"
                        href={"https://blockstream.info/"+(status.blockchainInfo.chain==="test"?"testnet/":"")+"block/"+block.hash}
                        rel="noopener norefferer"
                        target="_blank"
                      >
                        { block.hash.replace(/^0+/,'').substr(0,20 )+"..." }
                      </CLink>

                    </CCardHeader>
                    <CCardBody className="p-1">
                      <table className="table table-borderless p-0 m-0">
                        <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Height:</td><td className="p-0 pr-1 pl-1 m-0">{block.height}</td></tr>
                        <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Time:</td><td className="p-0 pr-1 pl-1 m-0">{new Date(block.time*1000).toLocaleString()}</td></tr>
                        <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Tx count:</td><td className="p-0 pr-1 pl-1 m-0">{block.nTx}</td></tr>
                        <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Size:</td><td className="p-0 pr-1 pl-1 m-0">{(block.size/1024/1024).toFixed(2) + " Mb"}</td></tr>
                      </table>
                    </CCardBody>
                  </CCard>
                </CTooltip>
                <CIcon name="cil-arrow-thick-right"  width="16"/>
                  {
                    status.latestBlocks?.length - 1 === index && (
                      <span className="font-weight-bold font-lg">&nbsp;...</span>
                    )
                  }
                </div>
                ))}
            </div>
          </CCardFooter>
      </CCard>
      <CCard>
        <CCardHeader>
          Intalled Apps
        </CCardHeader>
        <CCardBody>
          <div className="d-flex flex-row flex-wrap justify-content-center">
            {
              appList.data.map((appData) => (

                <CCard>
                  <CCardHeader className="d-flex justify-content-between align-items-center">
                    <div className="font-weight-bold">{appData.name}</div>
                    <CImg className="ml-0 mr-0 p-1" style={{ "background-color": appData.meta?.color, "border-radius":"10px"}} width="30" height="30" src={appData.meta?.icon}/>
                  </CCardHeader>
                  <CCardBody style={{minWidth: "250px"}}>
                    <table className="table-borderless flex-fill font-xs m-0">
                      <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Version:</td><td className="p-0 pr-1 pl-1 m-0">{appData.version}</td></tr>
                      <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Mount:</td><td className="p-0 pr-1 pl-1 m-0">/{appData.mountPoint}</td></tr>
                      { appData.container?.state && (
                        <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">State:</td><td className="p-0 pr-1 pl-1 m-0">{appData.container.state}</td></tr>
                      )}
                      { appData.container?.created && (
                        <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Created:</td><td className="p-0 pr-1 pl-1 m-0">{new Date(appData.container.created*1000).toLocaleString()}</td></tr>
                      )}
                    </table>
                  </CCardBody>
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
                </CCard>
              ))
            }
          </div>
        </CCardBody>
      </CCard>
    </>
  )
}

export default Dashboard
