import { Component, OnInit } from '@angular/core';
import { UsersService } from '../../services/users.service';
import { ServerDataSource } from 'ng2-smart-table';

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
      id: {
        title: 'ID',
        type: 'number',
      },
      login: {
        title: 'Login',
        type: 'string',
      },
      name: {
        title: 'Name',
        type: 'string',
      },
      emailAddress: {
        title: 'E-mail',
        type: 'string',
      },
    },
  };

  source: ServerDataSource;

  constructor(
    private usersService: UsersService,
  ) {
    this.source = usersService.getDataSource();
  }

  ngOnInit() {
  }

}
