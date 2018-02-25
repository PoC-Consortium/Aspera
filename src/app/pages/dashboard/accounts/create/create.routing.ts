import { Routes, RouterModule }  from '@angular/router';

import { CreateComponent } from './create.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
  {
    path: '',
    component: CreateComponent,
    children: [
        //{ path: 'seed', component: SeedComponent },
        //{ path: 'record', component: RecordComponent },
        //{ path: 'reproduce', component: ReproduceComponent },
        //{ path: 'verify', component: VerifyComponent },
    ]
  }
];

export const CreateRouting = RouterModule.forChild(routes);
