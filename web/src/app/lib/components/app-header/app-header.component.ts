import { Component, Input, OnInit } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { MarketService, StoreService, AccountService } from '../../services';
import { Account } from '../../model';
import { MatDialog } from '@angular/material';
import { SendBurstDialogComponent } from '../send-burst-dialog/send-burst-dialog.component';
import { HttpClient } from '@angular/common/http';


@Component({
    selector: 'app-header',
    styleUrls: ['./app-header.component.scss'],
    templateUrl: './app-header.component.html'
})
export class AppHeaderComponent implements OnInit {

    public selectedAccount: Account;
    public accounts: Account[];

    constructor(
        private marketService: MarketService,
        private storeService: StoreService,
        private accountService: AccountService,
        public dialog: MatDialog,
    ) {}

    ngOnInit() {
        this.accountService.currentAccount.subscribe((account) => {
            this.getSelectedAccounts();
            this.selectedAccount = account;
        })
    }

    private getSelectedAccounts() {
        this.storeService.getAllAccounts().then((accounts) => {
            this.accounts = accounts;
        });
    }

    public openSendDialog(): void {

        // get suggested fees
        this.marketService.getSuggestedFees().subscribe((fees) => {

            // open dialog
            const dialogRef = this.dialog.open(SendBurstDialogComponent, {
                width: '600px',
                data: { account: this.selectedAccount, fees: fees }
            });
    
            dialogRef.afterClosed().subscribe(result => {
                console.log('The dialog was closed', result);
            });
        });

    }

    public getTotalBurst(accounts: Account[]) {
        return accounts.reduce(((acc, { balance }) => acc + balance), 0)
    }

    public selectAccount(account: Account) {
        this.accountService.selectAccount(account);
    }

    public getPriceBTC(): string {
        const price = this.marketService.getCurrentBurstPriceBTC();
        if (price.length) {
            return `${price} BTC`;
        } else {
            return `Fetching price...`;
        }
    }

    public getTrendingIcon(): string {
        return Math.sign(parseFloat(this.marketService.getBurst24hChange())) === 1 ? 'trending_up' : 'trending_down';
    }
}
