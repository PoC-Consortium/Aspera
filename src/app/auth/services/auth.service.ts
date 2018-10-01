import { Injectable } from '@angular/core';
import { Observable, of, throwError } from 'rxjs';

import { Credentials } from '../models/credentials';
import { Account, Keys } from '../../lib/model';
import { AccountService, StoreService } from '../../lib/services';
import { fromPromise } from 'rxjs/internal-compatibility';
import { tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  constructor(private accountService: AccountService,
              private storeService: StoreService) {}

  login({ passphrase }: Credentials): Observable<Account> {
    /**
     * Simulate a failed login to display the error
     * message for the login form.
     */
    // if (passphrase !== 'test') {
    //   return throwError('Invalid passphrase');
    // }
    return fromPromise(this.accountService.createActiveAccount(passphrase)).pipe(
      tap((account) => {
        console.log(account);
        return fromPromise(this.storeService.selectAccount(account))
      }),
      tap((account) => {
        console.log(account);
        let syncAccount = new Account(account); // account is immutable
        return this.accountService.synchronizeAccount(syncAccount);
      })
    );  


    // getPublicKeyFromPassphrase
    // getAccountIdFromPublicKey
    // set account and accountRS
    // test account (main net): misery accept snow brave crazy avoid thank dwell itself still magic stretch
    // return of(new Account({
    //     id :"9137890273881363297",
    //     address: "BURST-66V3-TAFM-K2D7-9G4NX",
    //     keys: new Keys({ publicKey: "4bea9cc41dbb99ae0b6e594b9fcf4122e06023aa5bf9f831e275827282ff1f41" })
    // }))
  }

  logout() {
    return of(true);
  }
}
