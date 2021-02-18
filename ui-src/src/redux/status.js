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

import {__fetchContainerInfoByHash, __fetchStatus} from "../api";
import {fetchStatusError, fetchStatusPending, fetchStatusSuccess} from './actions';

function base64UrlEncode(str) {
  return btoa(str).replace(/=+$/,""); //.replace(/\+/g, '-').replace(/\//g, '_').replace(/\=+$/, '');
}

export function fetchStatus(session ) {
  return dispatch => {
    dispatch(fetchStatusPending());
    __fetchStatus(session)
      .then(res => {
        if( res.status !== 200 ) {
          throw( new Error("invalid fetchStatus code") );
        }
        return res.json();
      })
      .then(res => {
        if(res.error) {
          throw(res.error);
        }

        (new Promise( async (resolve, reject) => {

          for ( const feature of res.cyphernodeInfo.features.concat( res.cyphernodeInfo.optional_features ) ) {

            if ( !feature.active ) {
              feature.container = {
                state: "not running",
              };
              continue;
            }

            const base64Image = base64UrlEncode( feature.docker.ImageName+":"+feature.docker.Version );
            const containerResponse = await __fetchContainerInfoByHash( base64Image, session );

            if ( containerResponse && containerResponse.status === 200 ) {
              const json = await containerResponse.json();
              const networks = []
              Object.keys(json.NetworkSettings.Networks).map(network => {
                networks.push({
                  name: network,
                  ipAddress: json.NetworkSettings.Networks[network].IPAddress
                });
                return network;
              });
              feature.container = {
                id: json.Id,
                state: json.State,
                created: json.Created,
                networks: networks.sort( (a,b) => a.name.localeCompare(b.name) ),
                mounts: (json.Mounts||[]).map( mount => mount.Source ).filter( mount => !!mount ).sort( (a,b) => a.localeCompare(b) )
              };

            } else {
              feature.container = {
                state: "not running",
              };
            }

          }

          resolve(res);
        })).then( (status) => {
          dispatch(fetchStatusSuccess(status));
          return status;
        }).catch( error => {
          throw(error);
        })

      })
      .catch(error => {
        dispatch(fetchStatusError(error));
      })
  }
}
