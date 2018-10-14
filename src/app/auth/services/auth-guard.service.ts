import { Injectable } from '@angular/core';
import { CanActivate } from '@angular/router';
import { Store, select } from '@ngrx/store';
import { Observable, forkJoin } from 'rxjs';
import { map, take, concatMap, catchError, withLatestFrom, filter, first, concatAll, switchMap } from 'rxjs/operators';
import { AuthApiActions } from '../actions';
import * as fromAuth from '../reducers';
import { StoreService } from '../../lib/services';
import { fromPromise } from 'rxjs/internal-compatibility';
import { Account } from '../../lib/model';
import { AccountsComponent } from '../../pages/dashboard/accounts';
import { mergeMap } from 'rxjs-compat/operator/mergeMap';

@Injectable({
  providedIn: 'root',
})
export class AuthGuard implements CanActivate {
  authorized: boolean = false;
  constructor(private store: Store<fromAuth.State>,
              private storeService: StoreService) {
  }

  canActivate(): Observable<boolean> {
    return this.storeService.ready.pipe(
      filter(Boolean),
      switchMap(async (ready) => {
        const selectedAccount = await this.storeService.getSelectedAccount().catch(() => {});
        const allAccounts = await this.storeService.getAllAccounts().catch(() => {});
        console.log(selectedAccount, allAccounts);
        if (selectedAccount) {
          return true;
        } else if (allAccounts && allAccounts.length) {
          this.store.dispatch(new AuthApiActions.AccountsListRedirect());
          return false;
        } else {
          this.store.dispatch(new AuthApiActions.LoginRedirect());
          return false;
        }
      })
    )
  }
}
