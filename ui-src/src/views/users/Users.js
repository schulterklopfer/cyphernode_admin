import React, {useState, useEffect, useRef, useContext} from 'react'

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
import requests from "../../requests";

const Users = () => {

  const context = useContext( SessionContext );

  const [userList, setUserList] = useState({data: []})
  const [appList, setAppList] = useState({data: []})
  const [roleSuggestions, setRoleSuggestions] = useState([]);

  const [userToEdit, setUserToEdit] = React.useState({id: -1});
  const [userIdToDelete, setUserIdToDelete] = React.useState(0);
  const [appLookup, setAppLookup] = React.useState({});

  const loginRef = useRef("");
  const nameRef = useRef("");
  const emailAddressRef = useRef("");
  const passwordRef = useRef("");

  const [newRoles, setNewRoles] = React.useState([]);

  useEffect(() => {
    // component mount
    let ignore = false;

    async function fetchLists() {
      if (ignore) {
        return;
      }
      let appListResponse;

      try {
        appListResponse = await requests.getApps( context.session );
      } catch (e) {
        console.log( e );
      }
      const alu = {}
      if ( appListResponse && appListResponse.status === 200) {
        const al = appListResponse.body;
        const rs = [];
        for (const app of al.data) {
          alu[app.id] = app;
          for (const role of app.availableRoles) {
            rs.push({
              id: role.id,
              appId: role.appId,
              name: app.name+" - "+role.name
            });
          }
        }
        setRoleSuggestions(rs);
        setAppList(al);
        setAppLookup(alu);
      }

      let userListResponse;

      try {
        userListResponse = await requests.getUsers( context.session );
      } catch (e) {
        console.log( e );
      }

      if (userListResponse && userListResponse.status === 200) {
        // everything is ok
        const ul = userListResponse.body;
        for (const user of ul.data) {
          const localAppList = [];
          for (const role of user.roles) {
            let app = alu[role.appId];
            if (localAppList.indexOf(app) === -1) {
              localAppList.push(app);
            }
          }
          localAppList.sort((a, b) => {
            return a.name > b.name ? 1 : -1
          });
          user.appList = localAppList;
        }
        setUserList(ul);
      }
    }

    fetchLists();
    return () => {
      // component unmount
      ignore = true
    }
  }, [context.session])

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

    try {
      if (userToEdit.id === 0) {
        // create new user in backend
        const createData = {
          login: loginRef.current.value,
          name: nameRef.current.value,
          email_address: emailAddressRef.current.value,
          password: passwordRef.current.value,
          roles: newRoles
        };

        const userCreateResponse = await requests.createUser( createData, context.session );
        const patchedUser = userCreateResponse.body;
        userList.data.push( patchedUser );
        userList.data.sort( (a,b) => a.login>b.login?1:-1 );
        setUserList(userList);

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

        const userPatchResponse = await requests.patchUser( userToEdit.id, patchData, context.session );

        const patchedUser = userPatchResponse.body;
        console.log( patchedUser );
        const userIndex = userList.data.findIndex( user => user.id === patchedUser.id );
        if ( userIndex !== -1 ) {
          userList.data[userIndex]=patchedUser;
          setUserList(userList);
        }
      }
      clearUserToEdit();
    } catch (e) {
      console.log( e );
    }
  }

  const handleEditUserDelete = async (evt) => {
    evt.preventDefault();
    if (userToEdit.id < 0) {
      // what!?
      return;
    }

    if (userIdToDelete > 0) {
      try {
        const userDeleteResponse = await requests.deleteUser( userIdToDelete, context.session );

        if ( userDeleteResponse.status===204) {
          const userIndex = userList.data.findIndex( user => user.id === userIdToDelete );
          if ( userIndex !== -1 ) {
            userList.data.splice( userIndex, 1);
            setUserList(userList);
          }
        }
        clearUserToEdit();
        setUserIdToDelete(0);
        return;

      } catch (e) {
        console.log( e );
        return;
      }

    }
    setUserIdToDelete(userToEdit.id);
  }

  return (

    <>
      <CRow><CCol>
        <CCard>
          <CCardBody>
            User data: <pre>{ JSON.stringify(context.session.payload, null, 2) }</pre>
          </CCardBody>
        </CCard>
      </CCol></CRow>
      <CRow> {userList.data.map((userData) => (
        <UserDetails
          key={userData.id}
          userData={userData}
          appList={appList.data}
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
        show={userToEdit.id >= 0}
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
            {userToEdit.id > 0 && (
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
