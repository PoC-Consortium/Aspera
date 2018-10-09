import { Routes, RouterModule }  from '@angular/router';

import { AccountsComponent } from './accounts.component';
import { AccountNewComponent } from '../setup/account/account.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
  {
    path: 'accounts',
    component: AccountsComponent
  }
];

export const AccountsRouting = RouterModule.forChild(routes);
