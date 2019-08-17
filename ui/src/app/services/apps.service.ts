import { of as observableOf,  Observable } from 'rxjs';
import { Injectable } from '@angular/core';
import { App } from '../shared/app';
import { HttpClient } from '@angular/common/http';
import {ServerDataSource} from 'ng2-smart-table';

@Injectable({
  providedIn: 'root',
})

export class AppsService {
  apps: App[];
  dataSource: ServerDataSource;

  constructor(
    http: HttpClient,
  ) {
    this.apps = [
      {
        id: 1,
        hash: 'hash1',
        name: 'Cyphernode admin',
        description: 'Cyphernode admin panel',
        roles: [
          {
            id: 1,
            name: 'admin',
            description: 'admin role',
            autoAssign: false,
          },
        ],
      },
      {
        id: 2,
        hash: 'hash2',
        name: 'General stats',
        description: 'General stats',
        roles: [
          {
            id: 2,
            name: 'stats',
            description: 'admin role',
            autoAssign: true,
          },
        ],
      },
    ];
    this.dataSource = new ServerDataSource( http, {
      endPoint: 'assets/apps.json',
    });
  }

  getApps(): Observable<App[]> {
    return observableOf( this.apps );
  }

  addApp( appData: App ): Observable<App> {
    // validate with json scheme
    if ( this.apps.findIndex( user => user.hash === appData.hash ) === -1 ) {
      this.apps.push( appData );
    }
    return observableOf(appData);
  }

  removeApp( id: number ) {
    const index = this.apps.findIndex( app => app.id === id );
    if ( index !== -1 ) {
      this.apps.splice(index, 1);
    }
  }

  getAppByHash( hash: string ): Observable<App> {
    return observableOf(this.apps.find( app => app.hash === hash ));
  }

  getAppById( id: number ): Observable<App> {
    return observableOf(this.apps.find( app => app.id === id ));
  }

  getDataSource(): ServerDataSource {
    return this.dataSource;
  }

}
