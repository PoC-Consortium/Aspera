import { Component, OnInit } from '@angular/core';
import { StoreService, AccountService } from '../../../lib/services';
import { Account } from '../../../lib/model';
import { Store } from '@ngrx/store';
import * as fromAuth from '../../reducers';
import { AccountsListActions } from '../../actions';

@Component({
  selector: 'app-accounts-list',
  templateUrl: './accounts-list.component.html',
  styleUrls: ['./accounts-list.component.scss']
})
export class AccountsListComponent implements OnInit {

  accounts: Account[];

  constructor(
      private storeService: StoreService,
      private accountService: AccountService,
      private store: Store<fromAuth.State>
  ) {
      console.log('accounts list');

      this.storeService.ready.subscribe((ready) => {
        this.storeService.getAllAccounts().then((accounts) => {
          console.log(accounts);
            this.accounts = accounts;
        })
      });
  }

  selectAccount(account: Account) {
    this.store.dispatch(new AccountsListActions.SelectAccount({ account: account }));
  }

  ngOnInit() {}

}
