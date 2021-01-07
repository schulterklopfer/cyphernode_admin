import React from 'react'

const ContainerInfo = containerInfo => {
  return (
    <>
      { containerInfo.state && (
          <tr><td className="p-0 pr-1 m-0 font-weight-bold">State:</td><td className="p-0 pr-1 pl-1 m-0">{containerInfo.state}</td></tr>
      )}
      { containerInfo.created && (
        <tr><td className="p-0 pr-1 m-0 font-weight-bold">Created:</td><td className="p-0 pr-1 pl-1 m-0">{new Date(containerInfo.created*1000).toLocaleString()}</td></tr>
      )}
      { containerInfo.networks && (
        <tr><td className="p-0 pr-1 m-0 font-weight-bold align-top">Nets:</td><td className="p-0 pr-1 pl-1 m-0 align-top">{containerInfo.networks.map( network=>(
          <div>{network.name} (IP: {network.ipAddress})</div>
        ))}</td></tr>
      )}
    </>
  )
};

export default ContainerInfo;
