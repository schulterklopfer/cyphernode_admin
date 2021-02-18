import React, {useContext, useEffect} from 'react'
import {
  TheContent,
  TheSidebar,
  TheFooter,
  TheHeader
} from './index'
import SessionContext from "../sessionContext";
import {useDispatch} from "react-redux";
import {fetchApps} from "../redux/apps";
import {fetchStatus} from "../redux/status";
import {fetchUsers} from "../redux/users"

const TheLayout = () => {

  const context = useContext( SessionContext )
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(fetchApps( context.session ))
  },[dispatch,context.session])

  useEffect(() => {
    dispatch(fetchUsers( context.session ))
  },[dispatch,context.session])

  useEffect( () => {
    const doStuff = () => {
      dispatch(fetchStatus( context.session ));
    }
    doStuff();
    const interval = setInterval(doStuff,  10*1000);
    return () => {
      clearInterval(interval);
    }
  }, [dispatch,context.session] )

  return (
    <div className="c-app c-default-layout">
      <TheSidebar/>
      <div className="c-wrapper">
        <TheHeader/>
        <div className="c-body">
          <TheContent/>
        </div>
        <TheFooter/>
      </div>
    </div>
  )
}

export default TheLayout
