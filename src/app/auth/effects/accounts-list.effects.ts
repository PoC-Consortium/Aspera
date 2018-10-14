import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material';
import { Router } from '@angular/router';
import { Actions, Effect, ofType } from '@ngrx/effects';
import { map } from 'rxjs/operators';
import { AccountsListActions } from '../actions';
import { AuthService } from '../services/auth.service';
import { AccountsListActionTypes } from '../actions/accounts-list.actions';
import { StoreService } from '../../lib/services';

@Injectable()
export class AccountsListEffects {

  @Effect({ dispatch: false })
  selectAccount$ = this.actions$.pipe(
    ofType<AccountsListActions.SelectAccount>(
      AccountsListActionTypes.SelectAccount
    ),
    map(({ payload: { account }}) => {
        console.log('Account Selected: ', account);
        return this.storeService.selectAccount(account).then(() => {
            this.router.navigate(['/']);
        });
    }),
  );


  constructor(
    private actions$: Actions,
    private authService: AuthService,
    private router: Router,
    private dialog: MatDialog,
    private storeService: StoreService
  ) {}
}
