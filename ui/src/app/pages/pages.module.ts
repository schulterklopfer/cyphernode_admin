import { NgModule } from '@angular/core';

import { NbMenuModule } from '@nebular/theme';
import { ThemeModule } from '../@theme/theme.module';
import { DashboardModule } from './dashboard/dashboard.module';
import { PagesRoutingModule } from './pages-routing.module';
import { UserlistModule} from './userlist/userlist.module';
import { ApplistModule } from './applist/applist.module';

import { PagesComponent } from './pages.component';

@NgModule({
  imports: [
    PagesRoutingModule,
    NbMenuModule,
    ThemeModule,
    DashboardModule,
    UserlistModule,
    ApplistModule,
  ],
  declarations: [
    PagesComponent,
  ],
})
export class PagesModule {
}
