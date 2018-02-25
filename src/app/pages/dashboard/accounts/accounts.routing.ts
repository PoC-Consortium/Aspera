import { Routes, RouterModule }  from '@angular/router';

import { AccountsComponent } from './accounts.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
  {
    path: '',
    component: AccountsComponent,
    children: [
        { path: 'create', loadChildren: './create/create.module#AccountsCreateModule' }
    ]
  }
];

export const AccountsRouting = RouterModule.forChild(routes);
