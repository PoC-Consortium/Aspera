import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginPageComponent } from './containers/login-page.component';
import { AccountsListComponent } from './components/accounts-list/accounts-list.component';
import { AccountsListResolver } from './components/accounts-list/accounts-list.resolver';

const routes: Routes = [
  { path: 'login', component: LoginPageComponent },
  { path: 'accounts', component: AccountsListComponent, resolve: { accounts: AccountsListResolver } },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
  providers: [AccountsListResolver]
})
export class AuthRoutingModule {}
