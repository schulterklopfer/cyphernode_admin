import React, { useState, useEffect, useContext } from 'react'
import {
  CCard,
  CCardHeader,
  CCardBody,
  CButton
} from '@coreui/react'
import requests from "../../requests";
import ContainerBasicInfo from "../_common/ContainerBasicInfo";
import ContainerInfo from "../_common/ContainerInfo";
import SessionContext from "../../sessionContext";

const base64UrlEncode = (str) => {
  return btoa(str).replace(/=+$/,""); //.replace(/\+/g, '-').replace(/\//g, '_').replace(/\=+$/, '');
}

const Features = () => {
  const [status, setStatus] = useState({})
  const [containerDetails, setContainerContainerDetails] = useState("");
  const context = useContext( SessionContext )

  useEffect( () => {
    async function fetchStatus() {
      const response = await requests.getStatus( context.session );
      if ( response && response.status === 200 &&
        response.body.cyphernodeInfo &&
        response.body.cyphernodeInfo.features &&
        response.body.cyphernodeInfo.optional_features) {
        // everything is ok

        for ( const feature of response.body.cyphernodeInfo.features.concat( response.body.cyphernodeInfo.optional_features ) ) {

          if ( !feature.active ) {
            feature.container = {
              state: "not running",
            };
            continue;
          }

          const base64Image = base64UrlEncode( feature.docker.ImageName+":"+feature.docker.Version );
          const containerResponse = await requests.getDockerContainerByImageHash( base64Image, context.session );

          if ( containerResponse && containerResponse.status === 200 ) {
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
  }, [context.session] )

  return (
    <>
      <CCard>
        <CCardHeader className="h5">
          Features
        </CCardHeader>
        <CCardBody>
          <div className="d-flex flex-row flex-wrap justify-content-center">
            {
              status.cyphernodeInfo?.features.concat(status.cyphernodeInfo?.optional_features).sort(
                (fa, fb) => fb.active - fa.active || fa.name.localeCompare(fb.name)
              ).map((feature, index) => (
                <CCard key={index} style={{minWidth:"300px"}} className={feature.active?"mr-2":"mr-2 text-muted"}>
                  <CCardHeader className="h6 text-center" color={ feature.active?"success":""}>
                    <div className={feature.active?"text-light":""}>{feature.name}</div>
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

export default Features
