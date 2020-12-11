import React from 'react'
import {
  CButton,
  CCard,
  CCardBody,
  CCol,
  CContainer,
  CForm,
  CInput,
  CInputGroup,
  CInputGroupPrepend,
  CInputGroupText, CModal, CModalBody, CModalFooter, CModalHeader, CModalTitle,
  CRow
} from '@coreui/react'
import CIcon from '@coreui/icons-react'
import SessionContext from "../../../sessionContext";
import { ErrorStrings } from "../../../errors";

class Login extends React.Component {

  constructor( props ) {
    super( props );
    this.state ={"username":"", "password":""};
    this.handleChange = this.handleChange.bind(this);
  }

  handleChange( event ) {
    if ( !event.target.name ) {
      return;
    }
    const state = {};
    state[event.target.name]=event.target.value;
    this.setState( state )
    event.preventDefault();
  }

  render() {
    return (
      <SessionContext.Consumer>
        {(context) => (
          <div className="c-app c-default-layout flex-row align-items-center">
            <CContainer>
              <CRow className="justify-content-center">
                <CCol md="8">
                  <CCard className="p-4">
                    <CCardBody>
                      <CModal
                        show={!!this.state.error}
                        onClose={() => this.setState({error: undefined})}
                        color="danger"
                        size="sm"
                      >
                        <CModalHeader closeButton>
                          <CModalTitle>Login or password wrong</CModalTitle>
                        </CModalHeader>
                        <CModalBody>{ this.state.error }</CModalBody>
                        <CModalFooter>
                          <CButton color="primary" onClick={() => this.setState({error: undefined})}>Ok</CButton>
                        </CModalFooter>
                      </CModal>
                      <CForm>
                        <h1>Login</h1>
                        <p className="text-muted">Sign In to your account</p>
                        <CInputGroup className="mb-3">
                          <CInputGroupPrepend>
                            <CInputGroupText>
                              <CIcon name="cil-user" />
                            </CInputGroupText>
                          </CInputGroupPrepend>
                          <CInput type="text" name="username" placeholder={JSON.stringify(context)} autoComplete="username" value={this.state.username} onChange={this.handleChange}/>
                        </CInputGroup>
                        <CInputGroup className="mb-4">
                          <CInputGroupPrepend>
                            <CInputGroupText>
                              <CIcon name="cil-lock-locked" />
                            </CInputGroupText>
                          </CInputGroupPrepend>
                          <CInput type="password" name="password" placeholder="Password" autoComplete="current-password" value={this.state.password} onChange={this.handleChange}/>
                        </CInputGroup>
                        <CRow>
                          <CCol xs="6">
                            <CButton color="primary" className="px-4" onClick={ async (event)=>{
                              const errorCode = await context.login( this.state.username, this.state.password, this.props.history );
                              if ( !errorCode ) {
                                this.props.history.push("/");
                                console.log( "success");
                              } else {
                                this.setState({ error: ErrorStrings[errorCode] })
                              }
                            } }>Login</CButton>
                          </CCol>
                          {/*
                          <CCol xs="6" className="text-right">
                            <CButton color="link" className="px-0">Forgot password?</CButton>
                          </CCol>
                          */}
                        </CRow>
                      </CForm>
                    </CCardBody>
                  </CCard>
                </CCol>
              </CRow>
            </CContainer>
          </div>
        )}
      </SessionContext.Consumer>

    )
  }
}

export default Login;
