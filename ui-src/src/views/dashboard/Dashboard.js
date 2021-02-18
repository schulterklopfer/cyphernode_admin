import React, {useState} from 'react'
import {
  CBadge,
  CButton,
  CCard,
  CCardBody,
  CCardFooter,
  CCardHeader,
  CCol,
  CImg,
  CLink,
  CPopover,
  CProgress,
  CRow
} from '@coreui/react'
import CIcon from "@coreui/icons-react";
import ContainerBasicInfo from "../_common/ContainerBasicInfo";
import ContainerInfo from "../_common/ContainerInfo";
import {getApps, getStatus} from "../../redux/selectors";
import { useSelector} from "react-redux";

const Dashboard = () => {

  const [containerDetails, setContainerContainerDetails] = useState("");

  const apps = useSelector( getApps );
  const status = useSelector( getStatus );

  return (
    <>
      <CCard>
        <CCardHeader className="h5 d-flex flex-row justify-content-between">
          <div>{ (status.blockchainInfo?.initialblockdownload)?"Initial block download":"Sync fetchStatus" }</div>
          <div><CBadge className="ml-1" shape="pill" color="primary">{status.blockchainInfo?.chain}</CBadge></div>
        </CCardHeader>
        <CCardBody>
          <CRow className="text-center">
            <CCol md sm="12" className="mb-sm-2 mb-0">
              <strong>{status.blockchainInfo?.blocks} of {status.blockchainInfo?.headers} blocks</strong>
              <CProgress
                showPercentage={true}
                showValue={true}
                className="mt-2"
                precision={2}
                color={status.blockchainInfo?.initialblockdownload?"primary":"success"}
                stripped={status.blockchainInfo?.initialblockdownload}
                animated={status.blockchainInfo?.initialblockdownload}
                max={status.blockchainInfo?.headers}
                value={status.blockchainInfo?.blocks}
              />
            </CCol>
          </CRow>
        </CCardBody>
        <CCardFooter>
          <div className="container-fluid mb-sm-2 mb-0 text-center">
            {
              [
                { field: "blocks", title:"Block count", transform: input => input },
                { field: "headers", title:"Header count", transform: input => input },
                { field: "difficulty", title:"Difficulty", transform: input => input },
                { field: "pruned", title:"Pruned", transform: input => { return input?"Yes":"No" } },
              ].map( (item, index) => (
                <span key={index}>
                  <span className="font-weight-bold">{item.title}:</span>
                  <span> {item.transform((status.blockchainInfo||{})[item.field]) } </span>
                  { index !== 3 && (<span>&mdash; </span>) }
                </span>
              ))
            }
          </div>
        </CCardFooter>
      </CCard>
      <CCard>
        <CCardHeader className="h5">
          Latest blocks
        </CCardHeader>
        <CCardBody>
          <div className="d-flex flex-row flex-wrap justify-content-center">
            { status.latestBlocks?.map( (block, index) => (
              <div key={index} className="d-flex flex-row align-items-center">
                <CPopover
                  header="Hash"
                  content={ block.hash }
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
                        <tbody>
                          <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Height:</td><td className="p-0 pr-1 pl-1 m-0">{block.height}</td></tr>
                          <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Time:</td><td className="p-0 pr-1 pl-1 m-0">{new Date(block.time*1000).toLocaleString()}</td></tr>
                          <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Tx count:</td><td className="p-0 pr-1 pl-1 m-0">{block.nTx}</td></tr>
                          <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Size:</td><td className="p-0 pr-1 pl-1 m-0">{(block.size/1024/1024).toFixed(2) + " Mb"}</td></tr>
                        </tbody>
                      </table>
                    </CCardBody>
                  </CCard>
                </CPopover>
                <CIcon name="cil-arrow-thick-right"  width="16"/>
                {
                  status.latestBlocks?.length - 1 === index && (
                    <span className="font-weight-bold font-lg">&nbsp;...</span>
                  )
                }
              </div>
            ))}
          </div>
        </CCardBody>
      </CCard>
      <CCard>
        <CCardHeader className="h5">
          Intalled cypherapps
        </CCardHeader>
        <CCardBody>
          <div className="d-flex flex-row flex-wrap justify-content-center">
            {
              apps.map((appData, index) => (
                <CCard key={index} style={{minWidth: "300px"}} className="mr-2">
                  <CCardHeader className="h6 d-flex justify-content-between align-items-center">
                    <div>{appData.name}</div>
                    <CImg className="ml-0 mr-0 p-1" style={{ "backgroundColor": appData.meta?.color, "borderRadius":"10px"}} width="30" height="30" src={appData.meta?.icon}/>
                  </CCardHeader>
                  <CCardBody>
                    <table className="table-borderless flex-fill font-xs m-0 mb-4">
                      <tbody>
                        <tr><td className="p-0 pr-1 m-0 font-weight-bold">Version:</td><td className="p-0 pr-1 pl-1 m-0">{appData.version}</td></tr>
                        <tr><td className="p-0 pr-1 m-0 font-weight-bold">Mount:</td><td className="p-0 pr-1 pl-1 m-0">/{appData.mountPoint}</td></tr>
                        { appData.container && (
                          <ContainerBasicInfo {...appData.container} />
                        )}
                      </tbody>
                    </table>
                    {
                      (appData.container && (
                        <>
                          <ContainerInfo { ...{
                            containerInfo: appData.container,
                            meta: appData,
                            onClose: () => {
                              setContainerContainerDetails("");
                            },
                            show: containerDetails===appData.container.id
                          }}/>
                          <CButton
                            data-containerid={appData.container.id}
                            className="float-right"
                            size="sm"
                            shape="pill"
                            color="primary"
                            onClick={(event) => {
                              const containerId = event.target.dataset.containerid;
                              setContainerContainerDetails(containerId);
                            }}
                          >container info</CButton>
                        </>
                      ))
                    }
                  </CCardBody>
                  {
                    appData.id > 1?(
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
                    ):(<></>)
                  }
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

