import { Action } from '@ngrx/store';
import { Credentials } from '../../../../auth/models/credentials';

export enum AccountCreateActionsTypes {
  CreateActiveAccount = '[Create Page] Create Active Account',
  CreatePassiveAccount = '[Create Page] Create Passive Account',
}

export class CreateActiveAccount implements Action {
  readonly type = AccountCreateActionsTypes.CreateActiveAccount;

  constructor(public payload: Credentials) {}
}

export class CreatePassiveAccount implements Action {
  readonly type = AccountCreateActionsTypes.CreatePassiveAccount;

  constructor(public payload: object) {}
}

export type AccountCreateActionsUnion = CreateActiveAccount | CreatePassiveAccount;
