import { Component, OnInit } from '@angular/core';
import { Store, select } from '@ngrx/store';
import { Credentials } from '../models/credentials';
import * as fromAuth from '../reducers';
import { CreateActiveAccount } from '../../pages/dashboard/setup/account/create.actions';

@Component({
  selector: 'bc-login-page',
  template: `
    <h1>Burst</h1>
    <h2>Add an account to get started</h2>
    <app-account-create></app-account-create>
  `,
  styles: [],
})
export class LoginPageComponent implements OnInit {

  constructor(private store: Store<fromAuth.State>) {}

  ngOnInit() {}

  onSubmit(credentials: Credentials) {
    this.store.dispatch(new CreateActiveAccount(credentials));
  }
}
