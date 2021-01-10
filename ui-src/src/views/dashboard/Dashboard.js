import React, { useState, useEffect } from 'react'
import {
  CImg,
  CCard,
  CCardHeader,
  CCardBody,
  CCol,
  CProgress,
  CRow, CBadge, CCardFooter, CLink, CPopover, CButton
} from '@coreui/react'
import requests from "../../requests";
import CIcon from "@coreui/icons-react";
import ContainerBasicInfo from "./ContainerBasicInfo";
import ContainerInfo from "./ContainerInfo";

const base64UrlEncode = (str) => {
  return btoa(str).replace(/=+$/,""); //.replace(/\+/g, '-').replace(/\//g, '_').replace(/\=+$/, '');
}

const Dashboard = () => {
  const [status, setStatus] = useState({})
  const [appList, setAppList] = useState({ data: [] })
  const [containerDetails, setContainerContainerDetails] = useState("");

  useEffect( () => {
    async function fetchStatus() {
      const response = await requests.getStatus();
      if ( response.status === 200 ) {
        // everything is ok
        for ( const feature of response.body.cyphernodeInfo?.features?.concat( response.body.cyphernodeInfo?.optional_features ) ) {

          if ( !feature.active ) {
            feature.container = {
              state: "not running",
            };
            continue;
          }

          const base64Image = base64UrlEncode( feature.docker.ImageName+":"+feature.docker.Version );
          const containerResponse = await requests.getDockerContainerByImageHash( base64Image );


          if ( containerResponse.status === 200 ) {
            const networks = []
            Object.keys(containerResponse.body.NetworkSettings.Networks).map(network => {
              networks.push({
                name: network,
                ipAddress: containerResponse.body.NetworkSettings.Networks[network].IPAddress
              });
              return network;
            });
            feature.container = {
              id: containerResponse.body.Id,
              state: containerResponse.body.State,
              created: containerResponse.body.Created,
              networks: networks.sort( (a,b) => a.name.localeCompare(b.name) ),
              mounts: (containerResponse.body.Mounts||[]).map( mount => mount.Source ).filter( mount => !!mount ).sort( (a,b) => a.localeCompare(b) )
            };

          } else {
            feature.container = {
              state: "not running",
            };
          }

        }

        setStatus(response.body);
      }
    }

    fetchStatus();
    const interval = setInterval(fetchStatus,  10*1000);
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
          for ( const app of appList.data ) {
            const name = app.id === 1 ? 'cyphernodeadmin':app.hash;
            const containerResponse = await requests.getDockerContainerByName( name );

            if ( containerResponse.status === 200 ) {
              const networks = []
              Object.keys(containerResponse.body.NetworkSettings.Networks).map(network => {
                networks.push({
                  name: network,
                  ipAddress: containerResponse.body.NetworkSettings.Networks[network].IPAddress
                });
                return network;
              });
              app.container = {
                id: containerResponse.body.Id,
                state: containerResponse.body.State,
                created: containerResponse.body.Created,
                networks: networks.sort( (a,b) => a.name.localeCompare(b.name) ),
                mounts: (containerResponse.body.Mounts||[]).map( mount => mount.Source ).filter( mount => !!mount ).sort( (a,b) => a.localeCompare(b) )
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

      <CCard>
        <CCardHeader className="h5 d-flex flex-row justify-content-between">
          <div>{ (status.blockchainInfo?.initialblockdownload)?"Initial block download":"Sync status" }</div>
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
              appList.data.map((appData, index) => (
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
                        <ContainerInfo { ...{
                          containerInfo: appData.container,
                          meta: appData,
                          onClose: () => {
                            console.log( "modal closed" );
                            setContainerContainerDetails("");
                          },
                          show: containerDetails===appData.container.id
                        }}/>
                      ))
                    }
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
      <CCard>
        <CCardHeader className="h5">
          Features
        </CCardHeader>
        <CCardBody>
          <div className="d-flex flex-row flex-wrap justify-content-center">
            {
              status.cyphernodeInfo?.features.concat(status.cyphernodeInfo?.optional_features).sort(
                (fa, fb) => fb.active - fa.active || fa.label.localeCompare(fb.label)
              ).map((feature, index) => (

                <CCard key={index} style={{minWidth:"300px"}} className={feature.active?"mr-2":"mr-2 text-muted"}>
                  <CCardHeader className="h6 text-center" color={ feature.active?"success":""}>
                    <div className={feature.active?"text-light":""}>{feature.label}</div>
                  </CCardHeader>
                  <CCardBody className="font-xs" style={{minWidth: "250px"}}>
                    <table className="table-borderless flex-fill font-xs m-0 mb-4">
                      <tbody>
                        <tr><td className="p-0 pr-1 m-0 font-weight-bold">Image:</td><td className="p-0 pr-1 pl-1 m-0">{feature.docker?.ImageName}</td></tr>
                        <tr><td className="p-0 pr-1 m-0 font-weight-bold">Version:</td><td className="p-0 pr-1 pl-1 m-0">{feature.docker?.Version}</td></tr>
                        { feature.container && (
                          <ContainerBasicInfo {...feature.container} />
                        )}
                      </tbody>
                    </table>

                      {
                      (feature.container && (
                        <ContainerInfo { ...{
                          containerInfo: feature.container,
                          meta: {
                            name: feature.label,
                            version: feature.docker?.Version
                          },
                          onClose: () => {
                            console.log( "modal closed" );
                            setContainerContainerDetails("");
                          },
                          show: containerDetails===feature.container.id
                        }}/>
                      ))
                    }
                    { feature.active && (
                      <CButton
                      data-containerid={feature.container?.id}
                      className="float-right"
                      size="sm"
                      shape="pill"
                      color="primary"
                      onClick={(event) => {
                      const containerId = event.target.dataset.containerid;
                      setContainerContainerDetails(containerId);
                    }}
                      >container info</CButton>
                    ) }
                  </CCardBody>
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
