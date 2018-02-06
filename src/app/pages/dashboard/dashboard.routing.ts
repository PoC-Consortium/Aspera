import { Routes, RouterModule }  from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { ModuleWithProviders } from '@angular/core';

import { HomeComponent } from './home';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
    { path: 'create', loadChildren: '../create/create.module#CreateModule' },
    {
        path: 'dashboard',
        component: DashboardComponent,
        children: [
            { path: '', redirectTo: 'home', pathMatch: 'full' },
            { path: 'accounts', loadChildren: './accounts/accounts.module#AccountsModule' },
            { path: 'home', loadChildren: './home/home.module#HomeModule' },
            { path: 'transactions', loadChildren: './transactions/transactions.module#TransactionsModule' }
        ]
    }
];

export const routing: ModuleWithProviders = RouterModule.forChild(routes);
