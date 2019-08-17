import { Component, Input, OnInit } from '@angular/core';
import { Role } from '../../../shared/role';

import { ViewCell } from 'ng2-smart-table';

@Component({
  templateUrl: './roles-renderer.component.html',
  styleUrls: ['./roles-renderer.component.scss'],
})

export class RolesRendererComponent implements ViewCell, OnInit  {

  roles: Role[];

  @Input() value: any;
  @Input() rowData: any;

  ngOnInit() {
    this.roles = this.value;
  }

}
