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

import {__fetchApps, __fetchContainerInfoByName} from "../api";
import {fetchAppsError, fetchAppsPending, fetchAppsSuccess} from './actions';

export function fetchApps(session ) {
  return dispatch => {
    dispatch(fetchAppsPending());
    __fetchApps(session)
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
          for( const app of res.results ) {
            const name = app.id === 1 ? 'cyphernodeadmin':app.hash;
            try {
              const containerResponse = await __fetchContainerInfoByName( name );
              if ( containerResponse && containerResponse.status === 200 && containerResponse.body ) {
                const networks = []
                const json = await containerResponse.json();
                if( json.NetworkSettings ) {
                  Object.keys(json.NetworkSettings.Networks).map(network => {
                    networks.push({
                      name: network,
                      ipAddress: json.NetworkSettings.Networks[network].IPAddress
                    });
                    return network;
                  });
                }
                app.container = {
                  id: json.Id,
                  state: json.State,
                  created: json.Created,
                  networks: networks.sort( (a,b) => a.name.localeCompare(b.name) ),
                  mounts: (json.Mounts||[]).map( mount => mount.Source ).filter( mount => !!mount ).sort( (a,b) => a.localeCompare(b) )
                };
              }
            } catch (error) {
              reject(error);
            }
          }
          resolve(res.results);
        })).then( (apps) => {
          dispatch(fetchAppsSuccess(apps));
          return apps;
        }).catch( error => {
          throw(error);
        })

      })
      .catch(error => {
        dispatch(fetchAppsError(error));
      })
  }
}
