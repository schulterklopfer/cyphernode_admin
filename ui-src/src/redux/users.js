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
  fetchUsersError, fetchUsersPending, fetchUsersSuccess,
  createUserError, createUserPending, createUserSuccess,
  patchUserError, patchUserPending, patchUserSuccess,
  deleteUserError, deleteUserPending, deleteUserSuccess
} from "./actions";
import {__fetchUsers, __createUser, __patchUser, __deleteUser } from "../api";

export function fetchUsers(session ) {
  return dispatch => {
    dispatch(fetchUsersPending());
    __fetchUsers(session)
      .then(res => {
        if( res.status !== 200 ) {
          throw( new Error("invalid fetchUsers status code") );
        }
        return res.json();
      })
      .then(res => {
        if(res.error) {
          throw(res.error);
        }
        dispatch(fetchUsersSuccess(res.results));
      })
      .catch(error => {
        dispatch(fetchUsersError(error));
      })
  }
}

export function createUser(userData, session ) {
  return dispatch => {
    dispatch(createUserPending());
    __createUser(userData, session)
      .then(res => {
        if( res.status !== 201 ) { // status created
          throw( new Error("invalid createUser status code") );
        }
        return res.json();
      })
      .then(res => {
        dispatch(createUserSuccess(res));
      })
      .catch(error => {
        dispatch(createUserError(error));
      })
  }
}

export function patchUser( userId, userData, session ) {
  return dispatch => {
    dispatch(patchUserPending());
    __patchUser(userId, userData, session)
      .then(res => {
        if( res.status !== 200 ) {
          throw( new Error("invalid patchUser status code") );
        }
        return res.json();
      })
      .then(res => {
        dispatch(patchUserSuccess(res));
      })
      .catch(error => {
        dispatch(patchUserError(error));
      })
  }
}

export function deleteUser( userId, session ) {
  return dispatch => {
    dispatch(deleteUserPending());
    __deleteUser(userId, session)
      .then(res => {
        if( res.status !== 204 ) {
          throw( new Error("invalid deleteUser status code") );
        }
      })
      .then(() => {
        dispatch(deleteUserSuccess(userId));
      })
      .catch(error => {
        dispatch(deleteUserError(error));
      })
  }
}
