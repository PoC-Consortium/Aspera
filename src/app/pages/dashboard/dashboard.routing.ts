import { Routes, RouterModule }  from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { ModuleWithProviders } from '@angular/core';
import { 
    AuthGuardService as AuthGuard 
  } from '../../lib/services/auth-guard.service';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
    {
        path: 'dashboard',
        component: DashboardComponent,
        canActivate: [AuthGuard],
        children: [
            { path: '', redirectTo: 'home', pathMatch: 'full' },
            { path: 'accounts', loadChildren: './accounts/accounts.module#AccountsModule' },
            { path: 'accounts/create', loadChildren: './accounts/create/create.module#AccountsCreateModule' },
            { path: 'home', loadChildren: './home/home.module#HomeModule' },
            { path: 'setup', loadChildren: './setup/setup.module#SetupModule' },
            { path: 'transactions', loadChildren: './transactions/transactions.module#TransactionsModule' }
        ]
    }
];

export const routing: ModuleWithProviders = RouterModule.forChild(routes);
