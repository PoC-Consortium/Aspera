import { Routes, RouterModule }  from '@angular/router';

import { SetupComponent } from './setup.component';
import { AccountNewComponent } from './account/account.component';
import { CreateActiveAccountComponent } from './account/create-active/create.component';
import { CreatePassiveAccountComponent } from './account/create-passive/create-passive.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
    {
        path: 'setup',
        component: SetupComponent
    },
    {
        path: 'create',
        component: AccountNewComponent
    },
    {
        path: 'create/active',
        component: CreateActiveAccountComponent
    },
    {
        path: 'create/passive',
        component: CreatePassiveAccountComponent
    }
];

export const SetupRouting = RouterModule.forChild(routes);
