import { NgModule } from '@angular/core';
import { NbMenuModule } from '@nebular/theme';
import { NbCardModule } from '@nebular/theme';
import { Ng2SmartTableModule } from 'ng2-smart-table';
import { HttpClientModule } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { UserlistComponent } from './userlist.component';

@NgModule({
  declarations: [UserlistComponent],
  imports: [
    CommonModule,
    NbCardModule,
    NbMenuModule,
    NbCardModule,
    Ng2SmartTableModule,
    HttpClientModule,
  ],
})
export class UserlistModule { }
