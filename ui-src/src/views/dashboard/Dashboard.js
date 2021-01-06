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
                        href={"https://blockstream.info/testnet/block/"+block.hash}
                        rel="noopener norefferer"
                        target="_blank"
                      >
                        { block.hash.replace(/^0+/,'').substr(0,20 )+"..." }
                      </CLink>

                    </CCardHeader>
                    <CCardBody className="p-1">
                      <table className="table table-borderless p-0 m-0">
                        <tr><td className="p-0 pr-1 pl-1 m-0"><strong>Height:</strong></td><td className="p-0 pr-1 pl-1 m-0">{block.height}</td></tr>
                        <tr><td className="p-0 pr-1 pl-1 m-0"><strong>Time:</strong></td><td className="p-0 pr-1 pl-1 m-0">{new Date(block.time*1000).toLocaleString()}</td></tr>
                        <tr><td className="p-0 pr-1 pl-1 m-0"><strong>Tx count:</strong></td><td className="p-0 pr-1 pl-1 m-0">{block.nTx}</td></tr>
                        <tr><td className="p-0 pr-1 pl-1 m-0"><strong>Size:</strong></td><td className="p-0 pr-1 pl-1 m-0">{(block.size/1024/1024).toFixed(2) + " Mb"}</td></tr>
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
