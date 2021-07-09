import React, { useState } from 'react'
import {
  CCard,
  CCardHeader,
  CCardBody,
  CCardFooter, CLink, CPopover,
  CModal, CModalBody, CModalHeader, CModalTitle, CButton
} from '@coreui/react'
import {CopyToClipboard} from 'react-copy-to-clipboard';
import QRCode from 'qrcode.react'
import CIcon from "@coreui/icons-react";
import {useSelector} from "react-redux";
import {getStatus} from "../../redux/selectors";
import {resolveFile} from "../../api";

const onionValueTransforms = {
  "traefik": (feature) => { return {
    value: "http://"+feature.extra.tor_hostname+":"+(feature.extra.http_port||80),
    footer: (
      <CopyToClipboard text={"http://"+feature.extra.tor_hostname+":"+(feature.extra.http_port||80)}>
        <CButton size="sm" color="primary" className="py-0 my-0"><CIcon name="cil-clipboard" className="mr-1"/><span>to clipboard</span></CButton>
      </CopyToClipboard>
    )
  }},
  "bitcoin": (feature) => { return {
    value: "http://"+feature.extra.tor_hostname+":"+(feature.extra.port),
    footer: (
      <CopyToClipboard text={"http://"+feature.extra.tor_hostname+":"+feature.extra.port}>
        <CButton size="sm" color="primary" className="py-0 my-0"><CIcon name="cil-clipboard" className="mr-1"/><span>to clipboard</span></CButton>
      </CopyToClipboard>
    )
  }},
  "lightning": (feature) => { return {
    value: feature.extra.pubkey+"@"+feature.extra.tor_hostname+":"+feature.extra.port,
    footer: (
      <CopyToClipboard text={feature.extra.pubkey+"@"+feature.extra.tor_hostname+":"+feature.extra.port}>
        <CButton size="sm" color="primary" className="py-0 my-0"><CIcon name="cil-clipboard" className="mr-1"/><span>to clipboard</span></CButton>
      </CopyToClipboard>
    )
  }}

}

const onionValue = ( feature ) => {
  if ( onionValueTransforms[feature.label] ) {
    return onionValueTransforms[feature.label](feature);
  }
  return feature.extra.tor_hostname;
}

const Links = () => {
  const [qrZoom, setQrZoom] = useState({} );
  const status = useSelector( getStatus );

  const onions = []
  for ( const feature of status.cyphernodeInfo?.features?.concat( status.cyphernodeInfo?.optional_features ) ) {

    if (!feature.active) {
      feature.container = {
        state: "not running",
      };
      continue;
    }

    if ( feature.extra && feature.extra.torified ) {
      // tor is enable for traefik, show qr code
      // TODO: fix for o  === null
      const o = onionValue( feature );
      if( o ) {
        onions.push( {
          title: feature.name,
          value: o.value,
          footer: o.footer
        } );
      }
    }
  }

  return (
    <>
      <CCard>
        <CCardHeader className="h5">
          Important stuff
        </CCardHeader>
        <CCardBody>
          <div className="d-flex flex-row flex-wrap justify-content-center">
          {  onions.map( (onion, index) => (
            <CPopover
              key={index}
              content={ onion.value }
              placement="top"
            >
              <CCard className="mr-2 h-25">
                <CCardHeader className="h6 text-center">
                  {onion.title}
                </CCardHeader>
                <CCardBody>
                  <QRCode renderAs={"svg"} value={onion.value} onClick={()=>{
                    setQrZoom( onion )
                  }}/>
                </CCardBody>
                { onion.footer && (<CCardFooter className="d-flex flex-row flex-wrap justify-content-center">{onion.footer}</CCardFooter> ) }
              </CCard>
            </CPopover>) ) }
          </div>
          <div className="d-flex flex-row flex-wrap justify-content-center">
            <CCard className="mr-2">
              <CCardHeader className="h6 text-center">
                Config archive
              </CCardHeader>
              <CCardBody className="d-flex flex-row justify-content-center align-items-center">
                <CLink target="_blank" href={resolveFile("config.7z")}><CIcon name="cil-cloud-download" width={64}/></CLink>
              </CCardBody>
              <CCardFooter className="d-flex flex-row flex-wrap justify-content-center">
                config.7z
              </CCardFooter>
            </CCard>
            <CCard className="mr-2">
              <CCardHeader className="h6 text-center">
                Client archive
              </CCardHeader>
              <CCardBody className="d-flex flex-row justify-content-center align-items-center">
                <CLink target="_blank" href={resolveFile("client.7z")}><CIcon name="cil-cloud-download" width={64}/></CLink>
              </CCardBody>
              <CCardFooter className="d-flex flex-row flex-wrap justify-content-center">
                client.7z
              </CCardFooter>
            </CCard>
          </div>
        </CCardBody>
      </CCard>
      <CModal
        show={!!qrZoom.value}
        onClose={()=>setQrZoom({})}
        size={"lg"}
      >
        <CModalHeader>
          <CModalTitle>{qrZoom.title}</CModalTitle>
        </CModalHeader>
        <CModalBody className="text-center">
          <QRCode className="d-inline px-2 py-5" renderAs={"svg"} value={qrZoom.value||""} size={400} />
        </CModalBody>
      </CModal>
    </>
  )
}

export default Links
