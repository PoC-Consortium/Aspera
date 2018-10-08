import { Component, Input } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { MarketService, StoreService } from '../../services';
import { Account } from '../../model';
import { AccountsListActions } from '../../../auth/actions';
import { FormatInputPathObject } from 'path';
import { Store } from '@ngrx/store';
import * as fromAuth from '../../../reducers';


@Component({
    selector: 'app-header',
    styleUrls: ['./app-header.component.scss'],
    templateUrl: './app-header.component.html'
})
export class AppHeaderComponent {

    public selectedAccount: Account;
    public accounts: Account[];

    constructor(
        private marketService: MarketService,
        private storeService: StoreService,
        private store: Store<fromAuth.State>
    ) {
  
        this.storeService.ready.subscribe((ready) => {
            this.storeService.getAllAccounts().then((accounts) => {
                this.accounts = accounts;
            })
            this.storeService.getSelectedAccount().then((account) => {
                this.selectedAccount = account;
            })
          });

    }

    public getTotalBurst(accounts: Account[]) {
        return accounts.reduce(((acc, {balance}) => acc + balance), 0)
    }

    public selectAccount(account: Account) {
        this.store.dispatch(new AccountsListActions.SelectAccount({ account: account }));
    }
}
