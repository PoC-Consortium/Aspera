import { Component, OnInit, OnDestroy, ViewChild, Input } from '@angular/core';
import { Router } from '@angular/router';

import { MatTableDataSource, MatSort } from '@angular/material';
import { Transaction, Account } from '../../../lib/model';
import { AccountService, StoreService } from '../../../lib/services';
import { Converter } from '../../../lib/util';

@Component({
    selector: 'app-home',
    styleUrls: ['./home.component.css'],
    templateUrl: './home.component.html'
})
export class HomeComponent {
    displayedColumns = ['type', 'opposite', 'amount', 'fee', 'timestamp', 'confirmed'];
    recentTransactionData;
    account: Account;

    constructor(
        private storeService: StoreService,
        private accountService: AccountService
    ) {
        console.log('home');

        this.storeService.getSelectedAccount()
            .then((account) => {
                this.account = account;
                this.accountService.getTransactions(account.id)
                    .then((transactions) => {
                        console.log(transactions);
                        this.recentTransactionData = new MatTableDataSource<Transaction>(transactions);
                        this.recentTransactionData.sort = this.sort;
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

}