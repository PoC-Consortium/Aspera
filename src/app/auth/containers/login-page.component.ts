import { Component, OnInit } from '@angular/core';
import { Store, select } from '@ngrx/store';
import { Credentials } from '../models/credentials';
import * as fromAuth from '../reducers';
import { LoginPageActions } from '../actions';

@Component({
  selector: 'bc-login-page',
  template: `
    <app-account-create></app-account-create>
  `,
  styles: [],
})
export class LoginPageComponent implements OnInit {

  constructor(private store: Store<fromAuth.State>) {}

  ngOnInit() {}

  onSubmit(credentials: Credentials) {
    this.store.dispatch(new LoginPageActions.Login({ credentials }));
  }
}
