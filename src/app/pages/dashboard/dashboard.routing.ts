import { Routes, RouterModule }  from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { ModuleWithProviders } from '@angular/core';
import { LoginGuard } from '../login/login-guard.service';
import { AccountsComponent } from './accounts';
import { HomeComponent } from './home';
import { TransactionsComponent } from './transactions';

// noinspection TypeScriptValidateTypes
const routes: Routes = [
    {
        path: 'dashboard',
        component: DashboardComponent,
        canActivate: [LoginGuard],
        runGuardsAndResolvers: 'always',
        children: [
            { path: '', redirectTo: 'home', pathMatch: 'full' },
            { path: 'accounts', component: AccountsComponent },
            { path: 'home', component: HomeComponent },
            { path: 'setup', loadChildren: './setup/setup.module#SetupModule' },
            { path: 'transactions', component: TransactionsComponent }
        ]
    }
];

export const routing: ModuleWithProviders = RouterModule.forChild(routes);
