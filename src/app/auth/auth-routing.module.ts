import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginPageComponent } from './containers/login-page.component';
import { AccountsListComponent } from './components/accounts-list/accounts-list.component';

const routes: Routes = [
  { path: 'login', component: LoginPageComponent },
  { path: 'accounts', component: AccountsListComponent },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AuthRoutingModule {}
