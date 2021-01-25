import React, {useContext} from 'react'
import { useSelector, useDispatch } from 'react-redux'
import {
  CCreateElement,
  CSidebar,
  CSidebarBrand,
  CSidebarNav,
  CSidebarNavDivider,
  CSidebarNavTitle,
  CSidebarMinimizer,
  CSidebarNavDropdown,
  CSidebarNavItem,
} from '@coreui/react'

import CIcon from '@coreui/icons-react'

// sidebar nav config
import navigation from './_nav'
import SessionContext from "../sessionContext";

const showDebugNavItems = false;

const TheSidebar = () => {
  const dispatch = useDispatch();
  const show = useSelector(state => state.sidebarShow);
  const context = useContext( SessionContext );
  
  return (
    <CSidebar
      show={show}
      onShowChange={(val) => dispatch({type: 'set', sidebarShow: val })}
    >
      <CSidebarBrand className="d-md-down-none" to="/">
        <CIcon
          className="c-sidebar-brand-full"
          name="logo-negative"
          height={35}
        />
        <CIcon
          className="c-sidebar-brand-minimized"
          name="sygnet"
          height={35}
        />
      </CSidebarBrand>
      <CSidebarNav>

        <CCreateElement
          items={navigation.filter( n => {
            if ( !n.roles || !n.roles.length ) {
              return showDebugNavItems;
            }

            if ( !context.session ||
              !context.session.user ||
              !context.session.user.roles ||
              !context.session.user.roles.length ) {
              return false;
            }

            for( const rn of n.roles ) {
              const found = context.session.user.roles.findIndex( ro => rn === '*' || ro.name === rn );
              if ( found !== -1) {
                return true;
              }
            }
            return false;
          })}
          components={{
            CSidebarNavDivider,
            CSidebarNavDropdown,
            CSidebarNavItem,
            CSidebarNavTitle
          }}
        />
      </CSidebarNav>
      <CSidebarMinimizer className="c-d-md-down-none"/>
    </CSidebar>
  )
}

export default React.memo(TheSidebar)
