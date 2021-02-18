export const FETCH_APPS_PENDING = 'FETCH_APPS_PENDING';
export const FETCH_APPS_SUCCESS = 'FETCH_APPS_SUCCESS';
export const FETCH_APPS_ERROR = 'FETCH_APPS_ERROR';

export const FETCH_USERS_PENDING = 'FETCH_USERS_PENDING';
export const FETCH_USERS_SUCCESS = 'FETCH_USERS_SUCCESS';
export const FETCH_USERS_ERROR = 'FETCH_USERS_ERROR';

export const CREATE_USER_PENDING = 'CREATE_USER_PENDING';
export const CREATE_USER_SUCCESS = 'CREATE_USER_SUCCESS';
export const CREATE_USER_ERROR = 'CREATE_USER_ERROR';

export const PATCH_USER_PENDING = 'PATCH_USER_PENDING';
export const PATCH_USER_SUCCESS = 'PATCH_USER_SUCCESS';
export const PATCH_USER_ERROR = 'PATCH_USER_ERROR';

export const DELETE_USER_PENDING = 'DELETE_USER_PENDING';
export const DELETE_USER_SUCCESS = 'DELETE_USER_SUCCESS';
export const DELETE_USER_ERROR = 'DELETE_USER_ERROR';

export const FETCH_STATUS_PENDING = 'FETCH_STATUS_PENDING';
export const FETCH_STATUS_SUCCESS = 'FETCH_STATUS_SUCCESS';
export const FETCH_STATUS_ERROR = 'FETCH_STATUS_ERROR';

export const UI_SET_SIDEBAR_SHOW = 'UI_SET_SIDEBAR_SHOW';

export function setSidebarShow( sidebarShow ) {
  return {
    type: UI_SET_SIDEBAR_SHOW,
    sidebarShow: sidebarShow
  }
}

export function fetchAppsPending() {
  return {
    type: FETCH_APPS_PENDING
  }
}

export function fetchAppsSuccess(apps) {
  return {
    type: FETCH_APPS_SUCCESS,
    data: apps
  }
}

export function fetchAppsError(error) {
  return {
    type: FETCH_APPS_ERROR,
    error: error
  }
}

export function fetchUsersPending() {
  return {
    type: FETCH_USERS_PENDING
  }
}

export function fetchUsersSuccess(users) {
  return {
    type: FETCH_USERS_SUCCESS,
    data: users
  }
}

export function fetchUsersError(error) {
  return {
    type: FETCH_USERS_ERROR,
    error: error
  }
}

export function addUserPending() {
  return {
    type: CREATE_USER_PENDING
  }
}

export function addUserSuccess(user) {
  return {
    type: CREATE_USER_SUCCESS,
    data: user
  }
}

export function addUserError(error) {
  return {
    type: CREATE_USER_ERROR,
    error: error
  }
}

export function patchUserPending() {
  return {
    type: PATCH_USER_PENDING
  }
}

export function patchUserSuccess(user) {
  return {
    type: PATCH_USER_SUCCESS,
    data: user
  }
}

export function patchUserError(error) {
  return {
    type: PATCH_USER_ERROR,
    error: error
  }
}

export function deleteUserPending() {
  return {
    type: DELETE_USER_PENDING
  }
}

export function deleteUserSuccess(user) {
  return {
    type: DELETE_USER_SUCCESS,
    data: user
  }
}

export function deleteUserError(error) {
  return {
    type: DELETE_USER_ERROR,
    error: error
  }
}

export function fetchStatusPending() {
  return {
    type: FETCH_STATUS_PENDING
  }
}

export function fetchStatusSuccess(status) {
  return {
    type: FETCH_STATUS_SUCCESS,
    data: status
  }
}

export function fetchStatusError(error) {
  return {
    type: FETCH_STATUS_ERROR,
    error: error
  }
}
