import { Component, OnInit } from '@angular/core';
import { AppsService } from '../../services/apps.service';
import { ServerDataSource } from 'ng2-smart-table';
import { RolesRendererComponent } from './roles-renderer/roles-renderer.component';

@Component({
  selector: 'ngx-applist',
  templateUrl: './applist.component.html',
  styleUrls: ['./applist.component.scss'],
})
export class ApplistComponent implements OnInit {

  settings = {
    actions: {
      edit: false,
      add: false,
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
      hash: {
        title: 'Hash',
        type: 'string',
      },
      name: {
        title: 'Name',
        type: 'string',
      },
      description: {
        title: 'Description',
        type: 'string',
      },
      roles: {
        title: 'Available roles',
        type: 'custom',
        renderComponent: RolesRendererComponent,
      },
    },
  };

  source: ServerDataSource;

  constructor(
    private appsService: AppsService,
  ) {
    this.source = appsService.getDataSource();
  }

  ngOnInit() {
  }

}
