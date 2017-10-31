import { Routes, RouterModule }  from '@angular/router';

import { CreateComponent } from './create.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
    {
        path: '',
        component: CreateComponent,
    },
    { path: '**', redirectTo: '/accounts' }
];

export const routing = RouterModule.forChild(routes);
