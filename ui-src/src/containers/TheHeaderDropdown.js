import React, {useContext} from 'react'
import { useHistory } from "react-router-dom";

import {
  CDropdown,
  CDropdownItem,
  CDropdownMenu,
  CDropdownToggle,
} from '@coreui/react'
import CIcon from '@coreui/icons-react'
import SessionContext from "../sessionContext";

const TheHeaderDropdown = (props) => {
  const context = useContext( SessionContext )
  const history = useHistory();
  return (
    <CDropdown
      inNav
      className="c-header-nav-items mx-2"
      direction="down"
    >
      <CDropdownToggle className="c-header-nav-link" caret={false}>
        <CIcon name="cil-user"/>
      </CDropdownToggle>
      <CDropdownMenu className="pt-0" placement="bottom-end">
        <CDropdownItem
          header
          tag="div"
          color="light"
          className="text-center"
        >
          <strong>{context.session?.user?.name}</strong>
        </CDropdownItem>
        <CDropdownItem>
          <CIcon name="cil-heart" className="mfe-2" style={{'--ci-primary-color':'red'}}/>Profile
        </CDropdownItem>
        <CDropdownItem divider />
        <CDropdownItem onClick={()=>{history.push("/logout")}}>
          <CIcon name="cil-account-logout" className="mfe-2" />Logout
        </CDropdownItem>
      </CDropdownMenu>
    </CDropdown>
  )
}

export default TheHeaderDropdown
