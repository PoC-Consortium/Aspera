import { Routes, RouterModule }  from '@angular/router';

import { SetupComponent } from './setup.component';
import { AccountNewComponent } from './account/account.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
    {
        path: '',
        component: SetupComponent
    },
    {
        path: 'create',
        component: AccountNewComponent
    }
];

export const SetupRouting = RouterModule.forChild(routes);
