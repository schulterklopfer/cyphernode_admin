import React, { useState } from 'react'
import {
  CCard,
  CCardHeader,
  CCardBody,
  CButton
} from '@coreui/react'
import ContainerBasicInfo from "../_common/ContainerBasicInfo";
import ContainerInfo from "../_common/ContainerInfo";
import {useSelector} from "react-redux";
import {getStatus} from "../../redux/selectors";

const Features = () => {
  const [containerDetails, setContainerContainerDetails] = useState("");
  const status = useSelector( getStatus );

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
