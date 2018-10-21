import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginPageComponent } from './containers/login-page.component';
import { AccountsListComponent } from './components/accounts-list/accounts-list.component';
import { CreateActiveAccountComponent } from '../pages/dashboard/setup/account/create-active/create.component';
import { CreatePassiveAccountComponent } from '../pages/dashboard/setup/account/create-passive/create-passive.component';

const routes: Routes = [
  { path: 'welcome', component: LoginPageComponent },
  { path: 'welcome/active', component: CreateActiveAccountComponent },
  { path: 'welcome/passive', component: CreatePassiveAccountComponent },
  { path: 'accounts', component: AccountsListComponent },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AuthRoutingModule {}
