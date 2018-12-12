import { Routes, RouterModule }  from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { ModuleWithProviders } from '@angular/core';
import { LoginGuard } from '../login/login-guard.service';
import { AccountsComponent } from './accounts';
import { HomeComponent } from './home';
import { TransactionsComponent } from './transactions';
import { SetupComponent } from './setup/setup.component';
import { AccountNewComponent } from './setup/account/account.component';
import { CreateActiveAccountComponent } from './setup/account/create-active/create.component';
import { CreatePassiveAccountComponent } from './setup/account/create-passive/create-passive.component';
import { MessagesComponent } from './messages/messages.component';

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
            {
                path: 'setup',
                component: SetupComponent
            },
            {
                path: 'setup/create',
                component: AccountNewComponent
            },
            {
                path: 'setup/create/active',
                component: CreateActiveAccountComponent
            },
            {
                path: 'setup/create/passive',
                component: CreatePassiveAccountComponent
            },
            { path: 'transactions', component: TransactionsComponent },
            { path: 'messages', component: MessagesComponent }
        ]
    }
];

export const routing: ModuleWithProviders = RouterModule.forChild(routes);
