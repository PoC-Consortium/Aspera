import { Component, OnInit } from '@angular/core';
import { Store, select } from '@ngrx/store';
import { AccountService } from '../../lib/services';

@Component({
  selector: 'bc-login-page',
  template: `
    <h1>Burst</h1>
    <h2>Add an account to get started</h2>
    <app-account-new></app-account-new>
  `,
  styles: [],
})
export class LoginComponent implements OnInit {

  constructor(private accountService: AccountService) {}

  ngOnInit() {
    console.log('wtf');
  }

  onSubmit(credentials) {
    this.accountService.createActiveAccount(credentials);
  }
}
