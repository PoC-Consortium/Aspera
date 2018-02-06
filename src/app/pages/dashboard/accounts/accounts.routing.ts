import { Routes, RouterModule }  from '@angular/router';

import { AccountsComponent } from './accounts.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
  {
    path: '',
    component: AccountsComponent
  }
];

export const AccountsRouting = RouterModule.forChild(routes);
