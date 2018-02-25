import { Routes, RouterModule }  from '@angular/router';

import { SetupComponent } from './setup.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
    {
        path: '',
        component: SetupComponent,
    }
];

export const SetupRouting = RouterModule.forChild(routes);
