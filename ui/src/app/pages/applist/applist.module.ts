import { NgModule } from '@angular/core';
import { NbCardModule, NbIconModule, NbInputModule, NbTreeGridModule } from '@nebular/theme';
import { Ng2SmartTableModule } from 'ng2-smart-table';
import { HttpClientModule } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { ApplistComponent } from './applist.component';
import { RolesRendererComponent } from './roles-renderer/roles-renderer.component';

@NgModule({
  declarations: [ApplistComponent, RolesRendererComponent],
  entryComponents: [
    RolesRendererComponent,
  ],
  imports: [
    CommonModule,
    NbCardModule,
    NbTreeGridModule,
    NbIconModule,
    NbInputModule,
    Ng2SmartTableModule,
    HttpClientModule,
  ],
})
export class ApplistModule { }
