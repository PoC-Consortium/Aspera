import { Action } from '@ngrx/store';
import { Account } from '../../lib/model';

export enum AuthActionTypes {
  Logout = '[Auth] Logout',
  LogoutConfirmation = '[Auth] Logout Confirmation',
  LogoutConfirmationDismiss = '[Auth] Logout Confirmation Dismiss',
}

export class Logout implements Action {
  readonly type = AuthActionTypes.Logout;
  constructor(public payload: { account: Account }) {}
}

export class LogoutConfirmation implements Action {
  readonly type = AuthActionTypes.LogoutConfirmation;
  constructor(public payload: { account: Account }) {}
}

export class LogoutConfirmationDismiss implements Action {
  readonly type = AuthActionTypes.LogoutConfirmationDismiss;
}

export type AuthActionsUnion =
  | Logout
  | LogoutConfirmation
  | LogoutConfirmationDismiss;
