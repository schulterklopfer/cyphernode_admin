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

import {
  FETCH_APPS_ERROR,
  FETCH_APPS_PENDING,
  FETCH_APPS_SUCCESS,
  FETCH_USERS_ERROR,
  FETCH_USERS_PENDING,
  FETCH_USERS_SUCCESS,
  CREATE_USER_ERROR,
  CREATE_USER_PENDING,
  CREATE_USER_SUCCESS,
  PATCH_USER_ERROR,
  PATCH_USER_PENDING,
  PATCH_USER_SUCCESS,
  DELETE_USER_ERROR,
  DELETE_USER_PENDING,
  DELETE_USER_SUCCESS,
  FETCH_STATUS_ERROR,
  FETCH_STATUS_PENDING,
  FETCH_STATUS_SUCCESS,
  UI_SET_SIDEBAR_SHOW
} from './actions';

export const initialUiState = {
  sidebarShow: 'responsive'
}

export const initialAppsState = {
  pending: false,
  data: [],
  error: null
}

export const initialUsersState = {
  pending: false,
  data: [],
  error: null
}

export const initialStatusState = {
  pending: false,
  data: [],
  error: null
}

export function  uiReducer(state = initialUiState, { type, ...rest }) {
  switch (type) {
    case UI_SET_SIDEBAR_SHOW:
      return {...state, ...rest }
    default:
      return state
  }
}


export function appsReducer(state = initialAppsState, action) {
  switch(action.type) {
    case FETCH_APPS_PENDING:
      return {
        ...state,
        pending: true
      }
    case FETCH_APPS_SUCCESS:
      return {
        ...state,
        pending: false,
        data: action.data
      }
    case FETCH_APPS_ERROR:
      return {
        ...state,
        pending: false,
        error: action.error
      }
    default:
      return state;
  }
}

export function usersReducer(state = initialUsersState, action) {
  let userIndex = -1;
  let data;
  switch(action.type) {
    case FETCH_USERS_PENDING:
      return {
        ...state,
        pending: true
      }
    case FETCH_USERS_SUCCESS:
      return {
        ...state,
        pending: false,
        data: action.data
      }
    case FETCH_USERS_ERROR:
      return {
        ...state,
        pending: false,
        error: action.error
      }
    case CREATE_USER_PENDING:
      return {
        ...state,
        pending: true
      }
    case CREATE_USER_SUCCESS:
      data = state.data.concat(action.data);
      return {
        ...state,
        pending: false,
        data
      }
    case CREATE_USER_ERROR:
      return {
        ...state,
        pending: false,
        error: action.error
      }
    case PATCH_USER_PENDING:
      return {
        ...state,
        pending: true
      }
    case PATCH_USER_SUCCESS:
      // find index of user with same id in state.data:
      data = state.data.slice();
      userIndex = data.findIndex( user => user.id === action.data.id );
      // replace that user in state.data with action.data:
      if ( userIndex !== -1 ) {
        data[userIndex]=action.data;
      }
      return {
        ...state,
        pending: false,
        data
      }
    case PATCH_USER_ERROR:
      return {
        ...state,
        pending: false,
        error: action.error
      }
    case DELETE_USER_PENDING:
      return {
        ...state,
        pending: true
      }
    case DELETE_USER_SUCCESS:
      // find index of user with same id in state.data:
      data = state.data.slice();
      userIndex = data.findIndex( user => user.id === action.data );
      // replace that user in state.data with action.data:
      if ( userIndex !== -1 ) {
        data.splice( userIndex, 1);
      }
      return {
        ...state,
        pending: false,
        data
      }
    case DELETE_USER_ERROR:
      return {
        ...state,
        pending: false,
        error: action.error
      }
    default:
      return state;
  }
}


export function statusReducer(state = initialStatusState, action) {
  switch(action.type) {
    case FETCH_STATUS_PENDING:
      return {
        ...state,
        pending: true
      }
    case FETCH_STATUS_SUCCESS:
      return {
        ...state,
        pending: false,
        data: action.data
      }
    case FETCH_STATUS_ERROR:
      return {
        ...state,
        pending: false,
        error: action.error
      }
    default:
      return state;
  }
}
