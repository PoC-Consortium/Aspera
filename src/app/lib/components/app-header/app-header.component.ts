import { Component, Input } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { MarketService, StoreService } from '../../services';
import { Account } from '../../model';
import { AccountsListActions } from '../../../auth/actions';
import { FormatInputPathObject } from 'path';
import { Store } from '@ngrx/store';
import * as fromAuth from '../../../reducers';
import { MatDialog } from '@angular/material';
import { SendBurstDialogComponent } from '../send-burst-dialog/send-burst-dialog.component';
import { HttpClient } from '@angular/common/http';


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
        private store: Store<fromAuth.State>,
        public dialog: MatDialog,
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

    openSendDialog(): void {


        // get suggested fees
        this.marketService.getSuggestedFees().subscribe((fees) => {

            // open dialog
            const dialogRef = this.dialog.open(SendBurstDialogComponent, {
                width: '600px',
                data: { account: this.selectedAccount, fees: fees }
            });
    
            dialogRef.afterClosed().subscribe(result => {
                console.log('The dialog was closed', result);
                // this.store.dispatch(new AccountServiceActions.SendBurst({ result }))
            });
        });

    }

    public getTotalBurst(accounts: Account[]) {
        return accounts.reduce(((acc, { balance }) => acc + balance), 0)
    }

    public selectAccount(account: Account) {
        this.selectedAccount = account;
        this.store.dispatch(new AccountsListActions.SelectAccount({ account: account }));
    }
}
