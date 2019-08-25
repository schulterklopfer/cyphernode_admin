import { Component, OnInit } from '@angular/core';
import { UsersService } from '../../services/users.service';
import { RemoteDataSource } from '../../services/lib/remoteDataSource';

@Component({
  selector: 'ngx-userlist',
  templateUrl: './userlist.component.html',
  styleUrls: ['./userlist.component.scss'],
})
export class UserlistComponent implements OnInit {

  settings = {
    add: {
      addButtonContent: '<i class="nb-plus"></i>',
      createButtonContent: '<i class="nb-checkmark"></i>',
      cancelButtonContent: '<i class="nb-close"></i>',
    },
    edit: {
      editButtonContent: '<i class="nb-edit"></i>',
      saveButtonContent: '<i class="nb-checkmark"></i>',
      cancelButtonContent: '<i class="nb-close"></i>',
    },
    delete: {
      deleteButtonContent: '<i class="nb-trash"></i>',
      confirmDelete: true,
    },
    columns: {
      login: {
        title: 'Login',
        type: 'string',
      },
      password: {
        title: 'Password',
        type: 'password',
      },
      name: {
        title: 'Name',
        type: 'string',
      },
      email_address: {
        title: 'E-mail',
        type: 'string',
      },
    },
  };

  source: RemoteDataSource;

  constructor(
    private usersService: UsersService,
  ) {
    this.source = usersService.getDataSource();

    this.source.onAdded().subscribe( (a) => {

    } );
  }

  ngOnInit() {
  }

}
