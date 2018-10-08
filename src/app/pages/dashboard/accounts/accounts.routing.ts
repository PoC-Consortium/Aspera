import { Routes, RouterModule }  from '@angular/router';

import { AccountsComponent } from './accounts.component';
import { AccountNewComponent } from '../setup/account/account.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
  {
    path: '',
    component: AccountsComponent,
    children: [
        { path: 'create', component: AccountNewComponent }
    ]
  }
];

export const AccountsRouting = RouterModule.forChild(routes);
