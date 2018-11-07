import { Routes, RouterModule }  from '@angular/router';

import { TransactionsComponent } from './transactions.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
  {
    path: 'test',
    component: TransactionsComponent
  }
];

export const TransactionsRouting = RouterModule.forChild(routes);
