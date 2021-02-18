import React, {useContext} from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { setSidebarShow } from "../redux/actions";
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
import {getSidebarShow} from "../redux/selectors";

const showDebugNavItems = false;

const TheSidebar = () => {
  const dispatch = useDispatch();
  const show = useSelector( getSidebarShow );
  const context = useContext( SessionContext );
  const user = useSelector(state => {
    return state.users.data.find( user => user.id === context.session.jwt.id );
  });

  return (
    <CSidebar
      show={show}
      onShowChange={
        (val) => {
          dispatch(setSidebarShow(val))
        }
      }
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
              !user ||
              !user.roles ||
              !user.roles.length ) {
              return false;
            }

            for( const rn of n.roles ) {
              const found = user.roles.findIndex( ro => rn === '*' || ro.name === rn );
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
