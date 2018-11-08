import { Component, OnInit, OnDestroy, ViewChild, Input } from '@angular/core';
import { Router, NavigationEnd } from '@angular/router';

import { MatTableDataSource, MatSort } from '@angular/material';
import { Transaction, Account } from '../../../lib/model';
import { AccountService, StoreService } from '../../../lib/services';
import { Converter } from '../../../lib/util';
import { NotifierService } from 'angular-notifier';

@Component({
    selector: 'app-home',
    styleUrls: ['./home.component.css'],
    templateUrl: './home.component.html'
})
export class HomeComponent implements OnInit {
    displayedColumns = ['type', 'opposite', 'amount', 'fee', 'timestamp', 'confirmed'];
    recentTransactionData;
    account: Account;
    navigationSubscription;

    constructor(
        private storeService: StoreService,
        private accountService: AccountService,
        private router: Router,
        private notificationService: NotifierService
    ) {

        // handle route reloads (i.e. if user changes accounts)
        // this.navigationSubscription = this.router.events.subscribe((e: any) => {
        //     if (e instanceof NavigationEnd) {
        //         this.fetchTransactions();
        //     }
        // });
            
    }

    ngOnInit() {
        console.log('hi');
    }

    fetchTransactions() {
        this.storeService.getSelectedAccount()
            .then((account) => {
                this.account = account;
                this.accountService.getTransactions(account.id)
                    .then((transactions) => {
                        console.log(transactions);
                        this.recentTransactionData = new MatTableDataSource<Transaction>(transactions);
                        this.recentTransactionData.sort = this.sort;
                    },
                    (error) => {
                        // Todo: throw a warning to the user that their account is unverified!!
                        console.log(error);
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

    ngOnDestroy() {
        if (this.navigationSubscription) {  
            this.navigationSubscription.unsubscribe();
        }
    }

}