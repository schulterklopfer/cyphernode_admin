/*
 * MIT License
 *
 * Copyright (c) 2021 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

import React, {useContext, useRef} from 'react'
import {
  CBadge, CButton,
  CCard,
  CCardBody,
  CCardHeader,
  CForm,
  CFormGroup,
  CInput,
  CInputGroup,
  CInputGroupAppend, CInputGroupText,
  CLink, CModal,
  CModalBody, CModalFooter,
  CModalHeader,
  CModalTitle,
} from '@coreui/react'
import SessionContext from "../../sessionContext";
import {useDispatch, useSelector} from "react-redux";
import {getApps} from "../../redux/selectors";
import {CIcon} from "@coreui/icons-react";
import {patchUser} from "../../redux/users";

const UserProfile = () => {

  const context = useContext( SessionContext );
  const [editMode, setEditMode] = React.useState(false);
  const dispatch = useDispatch();

  const user = useSelector(state => {
    if ( !state || !state.users || !state.users.data || !state.users.data.length ) {
      return;
    }
    return state.users.data.find( user => user.id === context.session.jwt.id );
  });

  const nameRef = useRef("");
  const emailAddressRef = useRef("");
  const passwordRef = useRef("");

  const apps = useSelector( getApps );

  const handleEditUserSubmit = async (evt) => {
    evt.preventDefault();

    if (!editMode) {
      // what!?
      return;
    }

    if (
      nameRef.current.value === user.name &&
      emailAddressRef.current.value === user.email_address &&
      passwordRef.current.value === "" ) {
      setEditMode(false)
      return;
    }

    // patch user in backend if changed
    const patchData = {
    };

    if ( nameRef.current.value !== user.name ) {
      patchData.name = nameRef.current.value;
    }

    if ( emailAddressRef.current.value !== user.email_address ) {
      patchData.email_address = emailAddressRef.current.value;
    }

    if ( passwordRef.current.value !== "" ) {
      patchData.password = passwordRef.current.value;
    }

    dispatch( patchUser( user.id, patchData, context.session ) );

    setEditMode(false);

  }

  return (
    <>
    <CCard>
      <CCardHeader className="h5">
        Your profile
        <div className="card-header-actions">
          <CLink className="card-header-action" onClick={
            () => {
              nameRef.current.value = user.name;
              emailAddressRef.current.value = user.email_address;
              passwordRef.current.value = "";
              setEditMode(true);
            }
          }>
            <CIcon name="cil-pencil" />
          </CLink>
        </div>
      </CCardHeader>
      <CCardBody>
        <h5>General</h5>

        <table className="table table-borderless">
          <tbody>
          <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Login:</td><td className="p-0 pr-1 pl-1 m-0">{ user?.login }</td></tr>
          <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">name:</td><td className="p-0 pr-1 pl-1 m-0">{ user?.name }</td></tr>
          <tr><td className="p-0 pr-1 pl-1 m-0 font-weight-bold">Email address:</td><td className="p-0 pr-1 pl-1 m-0">{ user?.email_address }</td></tr>
          </tbody>
        </table>
        <h5>Roles</h5>
        { !!user.description  && <p  class="lead"><em>{ user.description}</em></p> }
        { apps.map( appData => {
          const filteredRoles = user.roles.filter( (r) => r.appId === appData.id );
          if ( filteredRoles.length > 0 ) {
            return (
              <CCard key={appData.id}>
                <CCardHeader>{appData.name}</CCardHeader>
                <CCardBody>
                  { filteredRoles.map( (role) => (<CBadge key={role.id} className="font-lg p-2 cna-role text-white mr-1" shape="pill">{role.name}</CBadge>) ) }
                </CCardBody>
              </CCard>
            );
          }
          return (<></>)
        } ) }
      </CCardBody>
    </CCard>
    <CModal
      show={ editMode }
      onClose={ () => setEditMode(false) }
      size="lg"
    >
      <CForm onSubmit={handleEditUserSubmit}>
        <CModalHeader closeButton>
          <CModalTitle>{user.login}</CModalTitle>
        </CModalHeader>
        <CModalBody>
          <CFormGroup>
            <CInputGroup>
              <CInput innerRef={nameRef} id="name" name="name"
                      placeholder="Name" autoComplete="name"/>
              <CInputGroupAppend>
                <CInputGroupText><CIcon name="cil-user"/></CInputGroupText>
              </CInputGroupAppend>
            </CInputGroup>
          </CFormGroup>
          <CFormGroup>
            <CInputGroup>
              <CInput innerRef={emailAddressRef} type="email" id="email"
                      name="email" placeholder="Email" autoComplete="email"/>
              <CInputGroupAppend>
                <CInputGroupText><CIcon name="cil-envelope-closed"/></CInputGroupText>
              </CInputGroupAppend>
            </CInputGroup>
          </CFormGroup>
          <CFormGroup>
            <CInputGroup>
              <CInput innerRef={passwordRef} type="password"
                      id="password" name="password" placeholder="Password"/>
              <CInputGroupAppend>
                <CInputGroupText><CIcon name="cil-asterisk"/></CInputGroupText>
              </CInputGroupAppend>
            </CInputGroup>
          </CFormGroup>
        </CModalBody>
        <CModalFooter>
          <CButton color="secondary" onClick={()=>setEditMode(false)}>Cancel</CButton>{' '}
          <CButton type="submit" color="primary">Save</CButton>
        </CModalFooter>
      </CForm>
    </CModal>
    </>
  )
}

export default UserProfile
