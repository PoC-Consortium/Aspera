import { Component, OnInit, OnDestroy, ViewChild, Input } from '@angular/core';
import { Router, NavigationEnd } from '@angular/router';

import { MatTableDataSource, MatSort, MatDialog } from '@angular/material';
import { Transaction, Account } from '../../../lib/model';
import { AccountService, StoreService, MarketService } from '../../../lib/services';
import { Converter, BurstUtil } from '../../../lib/util';
import { NotifierService } from 'angular-notifier';
import { SendBurstDialogComponent } from '../../../lib/components';

@Component({
    selector: 'app-home',
    styleUrls: ['./home.component.scss'],
    templateUrl: './home.component.html'
})
export class HomeComponent implements OnInit {
    displayedColumns = ['type', 'account', 'amount', 'fee', 'timestamp', 'confirmed'];
    recentTransactionData;
    account: Account;
    navigationSubscription;

    constructor(
        private storeService: StoreService,
        private accountService: AccountService,
        private router: Router,
        private notificationService: NotifierService,
        private marketService: MarketService,
        private dialog: MatDialog
    ) {

        // handle route reloads (i.e. if user changes accounts)
        this.navigationSubscription = this.router.events.subscribe((e: any) => {
            if (e instanceof NavigationEnd) {
                this.fetchTransactions();
            }
        });
            
    }

    convertFeeToBurst(feeNQT) {
        return BurstUtil.convertStringToNumber(feeNQT);
    }

    ngOnInit() {
        this.fetchTransactions();
    }

    fetchTransactions() {
        this.storeService.getSelectedAccount()
            .then((account) => {
                this.account = account;
                this.accountService.getTransactions(account.id)
                    .then((transactions) => {
                        this.recentTransactionData = new MatTableDataSource<Transaction>(transactions);
                        this.recentTransactionData.sort = this.sort;
                    },
                    (error) => {
                        // Todo: throw a warning to the user that their account is unverified!!
                        console.log(error);
                        this.recentTransactionData = [];
                        this.notificationService.notify('error', error.toString());
                    })
            })
    }

    applyFilter(filterValue: string) {
        filterValue = filterValue.trim(); // Remove whitespace
        filterValue = filterValue.toLowerCase(); // MatTableDataSource defaults to lowercase matches
        this.recentTransactionData.filter = filterValue;
    }

    @ViewChild(MatSort) sort: MatSort;

    public isOwnAccount(address: string): boolean {
        return address != undefined && address === this.account.address;
    }

    public convertTimestamp(timestamp: number): Date {
        return Converter.convertTimestampToDate(timestamp);
    }



    openSendDialog(): void {

        // get suggested fees
        this.marketService.getSuggestedFees().subscribe((fees) => {

                // open dialog
            const dialogRef = this.dialog.open(SendBurstDialogComponent, {
                width: '600px',
                data: { account: this.account, fees: fees }
            });

            dialogRef.afterClosed().subscribe(result => {
            console.log('The dialog was closed', result);
            });
        });
    }

    ngOnDestroy() {
        if (this.navigationSubscription) {  
            this.navigationSubscription.unsubscribe();
        }
    }

}