import { Routes, RouterModule }  from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { ModuleWithProviders } from '@angular/core';
import { AuthGuard } from '../../auth/services/auth-guard.service';
import { AccountNewComponent } from './setup/account/account.component';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
    {
        path: 'dashboard',
        component: DashboardComponent,
        canActivate: [AuthGuard],
        runGuardsAndResolvers: 'always',
        children: [
            { path: '', redirectTo: 'home', pathMatch: 'full' },
            { path: 'accounts', loadChildren: './accounts/accounts.module#AccountsModule' },
            { path: 'home', loadChildren: './home/home.module#HomeModule' },
            { path: 'setup', loadChildren: './setup/setup.module#SetupModule' },
            { path: 'transactions', loadChildren: './transactions/transactions.module#TransactionsModule' }
        ]
    }
];

export const routing: ModuleWithProviders = RouterModule.forChild(routes);
