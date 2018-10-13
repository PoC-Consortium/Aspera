import { Action } from '@ngrx/store';
import { Account } from '../../lib/model/account';

export enum AuthApiActionTypes {
  LoginSuccess = '[Auth/API] Login Success',
  LoginFailure = '[Auth/API] Login Failure',
  LoginRedirect = '[Auth/API] Login Redirect',
  AccountsListRedirect = '[Auth/API] Accounts List Redirect',
}

export class LoginSuccess implements Action {
  readonly type = AuthApiActionTypes.LoginSuccess;

  constructor(public payload: { account: Account }) {}
}

export class LoginFailure implements Action {
  readonly type = AuthApiActionTypes.LoginFailure;

  constructor(public payload: { error: any }) {}
}

export class LoginRedirect implements Action {
  readonly type = AuthApiActionTypes.LoginRedirect;
}

export class AccountsListRedirect implements Action {
  readonly type = AuthApiActionTypes.AccountsListRedirect;
}

export type AuthApiActionsUnion = LoginSuccess | LoginFailure | LoginRedirect | AccountsListRedirect;
