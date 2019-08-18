import { of as observableOf,  Observable } from 'rxjs';
import { Injectable } from '@angular/core';
import { User } from '../shared/user';
import { HttpClient } from '@angular/common/http';
import {ServerDataSource} from 'ng2-smart-table';

@Injectable({
  providedIn: 'root',
})
export class UsersService {
  users: User[];
  dataSource: ServerDataSource;

  constructor(
    http: HttpClient,
  ) {
    this.users = [
      {
        id: 1,
        login: 'skp',
        emailAddress: 'skp@skp.rocks',
        name: 'Markus',
      },
      {
        id: 2,
        login: 'kexkey',
        emailAddress: 'kexkey@kexkey.rocks',
        name: 'Etienne',
      },
    ];
    this.dataSource = new ServerDataSource( http, {
      endPoint: 'http://localhost:8080/api/v0/users',
    });
  }

  getUsers(): Observable<User[]> {
    return observableOf( this.users );
  }

  addUser( userData: User ): Observable<User> {
    // validate with json scheme
    if ( this.users.findIndex( user => user.login === userData.login ) === -1 ) {
      this.users.push( userData );
    }
    return observableOf(userData);
  }

  removeUser( id: number ) {
    const index = this.users.findIndex( user => user.id === id );
    if ( index !== -1 ) {
      this.users.splice(index, 1);
    }
  }

  getUserByLogin( login: string ): Observable<User> {
    return observableOf(this.users.find( user => user.login === login ));
  }

  getUserById( id: number ): Observable<User> {
    return observableOf(this.users.find( user => user.id === id ));
  }

  getDataSource(): ServerDataSource {
    return this.dataSource;
  }

}
