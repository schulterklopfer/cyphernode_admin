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

import React, {useRef, useContext} from 'react'

import {
  CCard,
  CCardBody,
  CRow,
  CCol,
  CModal,
  CForm,
  CModalHeader,
  CModalTitle,
  CModalBody,
  CFormGroup,
  CInputGroup,
  CInput,
  CInputGroupAppend,
  CInputGroupText, CCardHeader, CModalFooter, CButton
} from '@coreui/react'
import UserDetails from "./UserDetails";
import {CIcon} from "@coreui/icons-react";
import ReactTags from "react-tag-autocomplete";
import rolesAreEqual from "./rolesAreEqual";
import SessionContext from "../../sessionContext";
import {getApps, getUsers} from "../../redux/selectors";
import {useDispatch, useSelector} from "react-redux";
import {createUser, deleteUser, patchUser} from "../../redux/users";

const Users = () => {

  const dispatch = useDispatch();
  const context = useContext( SessionContext );

  const [userToEdit, setUserToEdit] = React.useState({id: -1});
  const [userIdToDelete, setUserIdToDelete] = React.useState(0);

  const loginRef = useRef("");
  const nameRef = useRef("");
  const emailAddressRef = useRef("");
  const passwordRef = useRef("");

  const [newRoles, setNewRoles] = React.useState([]);

  const apps = useSelector( getApps );
  const users = useSelector( getUsers );

  const appLookup = {}
  const roleSuggestions = [];
  for (const app of apps) {
    appLookup[app.id] = app;
    for (const role of app.availableRoles) {
      roleSuggestions.push({
        id: role.id,
        appId: role.appId,
        name: app.name+" - "+role.name
      });
    }
  }

  const appsByUserId = {};

  for (const user of users) {
    const localAppList = [];
    for (const role of user.roles) {
      let app = appLookup[role.appId];
      if (localAppList.indexOf(app) === -1) {
        localAppList.push(app);
      }
    }
    localAppList.sort((a, b) => {
      return a.name > b.name ? 1 : -1
    });
    appsByUserId[user.id]=localAppList;
  }

  const clearUserToEdit = () => {
    setUserIdToDelete(0);
    setUserToEdit({id: -1});
  }

  const handleEditUserSubmit = async (evt) => {
    evt.preventDefault();

    if (userToEdit.id === -1) {
      // what!?
      return;
    }

    if (
      loginRef.current.value === userToEdit.login &&
      nameRef.current.value === userToEdit.name &&
      emailAddressRef.current.value === userToEdit.email_address &&
      passwordRef.current.value === "" &&
      rolesAreEqual(newRoles, userToEdit.roles)) {
      clearUserToEdit();
      return;
    }

    if (userToEdit.id === 0) {
      // create new user in backend
      dispatch( createUser({
        login: loginRef.current.value,
        name: nameRef.current.value,
        email_address: emailAddressRef.current.value,
        password: passwordRef.current.value,
        roles: newRoles
      }, context.session ));

    } else {
      // patch user in backend if changed
      const patchData = {
      };

      if ( loginRef.current.value !== userToEdit.login ) {
        patchData.login = loginRef.current.value;
      }

      if ( nameRef.current.value !== userToEdit.name ) {
        patchData.name = nameRef.current.value;
      }

      if ( emailAddressRef.current.value !== userToEdit.email_address ) {
        patchData.email_address = emailAddressRef.current.value;
      }

      if ( passwordRef.current.value !== "" ) {
        patchData.password = passwordRef.current.value;
      }

      if ( !rolesAreEqual(newRoles, userToEdit.roles)) {
        // check which roles are new and which ones
        // we need to delete
        patchData.roles = newRoles;
      }

      dispatch( patchUser( userToEdit.id, patchData, context.session ) );
    }
    clearUserToEdit();

  }

  const handleEditUserDelete = async (evt) => {
    evt.preventDefault();
    if (userToEdit.id < 0) {
      // what!?
      return;
    }

    if (userIdToDelete > 0) {
      dispatch( deleteUser( userIdToDelete, context.session ) );
      clearUserToEdit();
      setUserIdToDelete(0);
      return;
    }
    setUserIdToDelete(userToEdit.id);
  }

  return (

    <>
      <CRow><CCol>
        <CCard>
          <CCardBody>
            User data: <pre>{ JSON.stringify(context.session.user, null, 2) }</pre>
          </CCardBody>
        </CCard>
      </CCol></CRow>
      <CRow> {users.map((userData) => (
        <UserDetails
          key={userData.id}
          userData={userData}
          appList={apps}
          roleSuggestions={roleSuggestions}
          onEditClick={(userData) => {
            loginRef.current.value = userData.login;
            nameRef.current.value = userData.name;
            emailAddressRef.current.value = userData.email_address;
            passwordRef.current.value = "";
            setNewRoles(userData.roles.slice() || []);
            setUserToEdit(userData);
          }}
        />))}
      </CRow>
      <CRow>
        <CCol className="mb-3 mb-xl-0">
          <CButton shape="pill" color="primary" className="float-right" onClick={
            () => {
              loginRef.current.value = "";
              nameRef.current.value = "";
              emailAddressRef.current.value = "";
              passwordRef.current.value = "";
              setNewRoles( []);
              setUserToEdit({
                id: 0, login: "", name: "", email_address: "", password: "", roles: []
              });
            }
          }><CIcon name="cil-user-plus" /> Add user</CButton>
        </CCol>
      </CRow>
      <CModal
        show={ userToEdit.id>=0 && userToEdit.id !== 1 }
        onClose={clearUserToEdit}
        color={userIdToDelete > 0 ? "danger" : "default"}
        size="lg"
      >
        <CForm onSubmit={handleEditUserSubmit}>
          <CModalHeader closeButton>
            <CModalTitle>{userIdToDelete === 0 ? userToEdit.id===0?"Create new user":"Edit details of" : "Delete user"} {userToEdit.login}</CModalTitle>
          </CModalHeader>
          <CModalBody>
            <fieldset disabled={userIdToDelete > 0 ? "disabled" : ""}>
              <CFormGroup>
                <CInputGroup>
                  <CInput innerRef={loginRef} id="login" name="login"
                          placeholder="Login" autoComplete="login"/>
                  <CInputGroupAppend>
                    <CInputGroupText><CIcon name="cil-at"/></CInputGroupText>
                  </CInputGroupAppend>
                </CInputGroup>
              </CFormGroup>
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

              <CCard>
                <CCardHeader>Roles</CCardHeader>
                <CCardBody>
                  <div style={{width: "100%"}}>
                    <ReactTags
                      foo={"bar"}
                      removeButtonText="Click to remove role"
                      tagComponent={(props) => (
                        React.createElement( 'button', { type: 'button', className: props.classNames.selectedTag, title: props.removeButtonText, onClick: props.onDelete },
                          React.createElement( 'span', { className: props.classNames.selectedTagName }, ((appLookup[props.tag.appId] || {}).name || "") + " - "  + props.tag.name )
                        )
                      )}
                      onAddition={(e) => {
                        if (newRoles.findIndex((role) => role.id === e.id) === -1) {
                          const app = appLookup[e.appId];
                          if( app.availableRoles && app.availableRoles.length ) {
                            const role = app.availableRoles.find( (ar) => ar.id === e.id );
                            if ( role ) {
                              newRoles.push(role);
                            }
                          }
                          newRoles.sort((a, b) => a.name > b.name ? 1 : -1);
                          setNewRoles(newRoles);
                        }
                      }}
                      onDelete={(index) => {
                        newRoles.splice(index, 1);
                        setNewRoles(newRoles);
                      }}
                      placeholderText="Add new role"
                      allowBackspace={false}
                      tags={newRoles}
                      suggestions={roleSuggestions}
                    />
                  </div>
                </CCardBody>
              </CCard>
            </fieldset>
          </CModalBody>
          <CModalFooter>
            {userToEdit?.id > 1 && (
              <>
                <CButton color="danger"
                         onClick={handleEditUserDelete}>{userIdToDelete > 0 ? "FINISH IT" : "Delete"}</CButton>
                <div className="flex-fill"/>
              </>
            )}
            <CButton color="secondary" onClick={clearUserToEdit}>Cancel</CButton>{' '}
            {userIdToDelete === 0 && (
              <CButton type="submit" color="primary">Save</CButton>
            )}
          </CModalFooter>
        </CForm>
      </CModal>

    </>

  )
}

export default Users
