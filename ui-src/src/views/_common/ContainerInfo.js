import React, { useRef, useState } from 'react'
import {CCardBody, CCard, CCardHeader, CModal, CModalHeader, CModalFooter, CModalBody, CRow, CCol} from "@coreui/react";
import requests from "../../requests";


const ContainerInfo = props => {

  const containerInfo = props.containerInfo;
  const meta = props.meta;

  const [logLines, setLogLines] = useState( [] );
  const [client, setClient] = useState(null)

  const endRef = useRef();

  const onModalOpened = () => {
    console.log( "on opened");
    const client = requests.getDockerLogsWebsocketClient( containerInfo.id )
    client.onopen = () => {
      console.log('WebSocket Client Connected');
    };
    client.onmessage = (message) => {
      setLogLines( logLines => {
        const newLines = JSON.parse( message.data );
        const newLogLines = [...logLines, ...newLines ];
        if( newLogLines.length > 512 ) {
          return newLogLines.slice( newLogLines.length - 512 );
        }
        return newLogLines;
      }  );
      endRef.current.scrollIntoView({ behavior: 'smooth' })
    };
    client.onerror = () => {
      client.close(1000);
      setClient(null);
    }
    setClient(client)
  }

  const onModalClosed = () => {
    console.log( "on closed");
    client.close(1000);
    setClient(null);
  }


  return (

    <CModal className="mw-100 w-100"
            show={props.show}
            onClose={props.onClose}
            onClosed={onModalClosed}
            onOpened={onModalOpened}
    >
      <CModalHeader className="h5" closeButton>{meta.name} {meta.version}</CModalHeader>
      <CModalBody>
        <CRow>
        {
          (containerInfo.networks && containerInfo.networks.length)?(
          <CCol>
            <CCard className="font-xs">
            <CCardHeader className="h6">Container networks</CCardHeader>
            <CCardBody>
            <pre className="m-0 p-0">
            {
              containerInfo.networks.map( (network, index) =>
                `[${index}] ${network.name} (${network.ipAddress})\n`
              )
            }
            </pre>
              </CCardBody>
            </CCard>
          </CCol>

          ):(<></>)
        }

        {
          (containerInfo.mounts && containerInfo.mounts.length)?(
          <CCol className="col-8">
            <CCard className="font-xs">
              <CCardHeader className="h6">Container mounts</CCardHeader>
              <CCardBody>
            <pre className="m-0 p-0">
            {
              containerInfo.mounts.map( (mount, index) =>
                `[${index}] ${mount}\n`
              )
            }
            </pre>
              </CCardBody>
            </CCard>
          </CCol>
          ):(<></>)
        }
        </CRow>
        <CRow>
          <CCol>
        {
          logLines.length?(
            <CCard className="font-xs">
              <CCardHeader className="h6">Container logs</CCardHeader>
              <CCardBody className="docker-logs">
              <pre className="m-0 p-0" style={{
                height: "300px"
              }}>
              {
                logLines.map( line => line+"\n" )
              }
                <div ref={endRef} />
              </pre>
              </CCardBody>
            </CCard>
          ):(<></>)

        }
          </CCol>
        </CRow>
      </CModalBody>
      <CModalFooter>
          <span><span className="font-weight-bold">State:&nbsp;</span> <span>{containerInfo.state}</span></span>
          <span><span className="font-weight-bold">Created:&nbsp;</span> <span>{new Date(containerInfo.created*1000).toLocaleString()}</span></span>
          <span><span className="font-weight-bold">ID:&nbsp;</span> <span>{containerInfo.id}</span></span>
      </CModalFooter>
    </CModal>
  )
};

export default ContainerInfo;
