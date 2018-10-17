import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material';
import { Router } from '@angular/router';
import { Actions, Effect, ofType } from '@ngrx/effects';
import { of, combineLatest, merge, pipe, forkJoin, concat } from 'rxjs'; 
import { catchError, exhaustMap, map, tap, mergeMap, mapTo, withLatestFrom, filter } from 'rxjs/operators';
import {
  LoginPageActions,
  AuthActions,
  AuthApiActions,
} from '../actions';
import { Credentials } from '../models/credentials';
import { LogoutConfirmationDialogComponent } from '../components/logout-confirmation-dialog.component';
import { Account } from '../../lib/model';
import { AccountService } from '../../lib/services';

@Injectable()
export class AuthEffects {
  @Effect()
  login$ = this.actions$.pipe(
    ofType<LoginPageActions.Login>(LoginPageActions.LoginPageActionTypes.Login),
    map(action => action.payload.credentials),
    exhaustMap((auth: Credentials) =>
      this.accountService.login(auth).pipe(
        map(account => new AuthApiActions.LoginSuccess({ account })),
        catchError(error => of(new AuthApiActions.LoginFailure({ error })))
      )
    )
  );

  @Effect({ dispatch: false })
  loginSuccess$ = this.actions$.pipe(
    ofType(AuthApiActions.AuthApiActionTypes.LoginSuccess),
    tap(() => this.router.navigate(['/']))
  );

  @Effect({ dispatch: false })
  loginRedirect$ = this.actions$.pipe(
    ofType(
      AuthApiActions.AuthApiActionTypes.LoginRedirect,
    ),
    tap(authed => {
      this.router.navigate(['/login']);
    })
  );

  @Effect({ dispatch: false })
  accountsListRedirect$ = this.actions$.pipe(
    ofType(
      AuthApiActions.AuthApiActionTypes.AccountsListRedirect,
    ),
    tap(authed => {
      this.router.navigate(['/accounts']);
    })
  );

  @Effect()
  logoutConfirmation$ = this.actions$.pipe(
    ofType<AuthActions.LogoutConfirmation>(AuthActions.AuthActionTypes.LogoutConfirmation),
    map(action => action.payload.account),
    exhaustMap((account: Account) => {
      const dialogRef = this.dialog.open<
        LogoutConfirmationDialogComponent,
        undefined,
        boolean
      >(LogoutConfirmationDialogComponent);

      return dialogRef.afterClosed().pipe(
        map((logout) => {
          return logout
            ? new AuthActions.Logout({ account: account })
            : new AuthActions.LogoutConfirmationDismiss()
        })
      );
    })
  );

  @Effect({ dispatch:false })
  logoutUser$ = this.actions$.pipe(
    ofType(
      AuthActions.AuthActionTypes.Logout
    ),
    tap(({ payload: { account} }: AuthActions.Logout) => {
      this.accountService.removeAccount(account);
      this.router.navigate(['/']);
    })
  );


  constructor(
    private actions$: Actions,
    private router: Router,
    private dialog: MatDialog,
    private accountService: AccountService
  ) {}
}
