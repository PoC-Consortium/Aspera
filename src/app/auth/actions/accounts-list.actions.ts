import { Action } from '@ngrx/store';
import { Account } from '../../lib/model/account';

export enum AccountsListActionTypes {
  SelectAccount = '[Accounts List] Select Account',
}

export class SelectAccount implements Action {
  readonly type = AccountsListActionTypes.SelectAccount;

  constructor(public payload: { account: Account }) {}
}


export type AccountsListActionsUnion = SelectAccount;
