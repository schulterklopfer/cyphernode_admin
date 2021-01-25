import React from 'react'

const ContainerBasicInfo = containerInfo => {
  return (
    <>
      { containerInfo.state && (
          <tr><td className="p-0 pr-1 m-0 font-weight-bold">State:</td><td className="p-0 pr-1 pl-1 m-0">{containerInfo.state}</td></tr>
      )}
      { containerInfo.created && (
        <tr><td className="p-0 pr-1 m-0 font-weight-bold">Created:</td><td className="p-0 pr-1 pl-1 m-0">{new Date(containerInfo.created*1000).toLocaleString()}</td></tr>
      )}
    </>
  )
};

export default ContainerBasicInfo;
