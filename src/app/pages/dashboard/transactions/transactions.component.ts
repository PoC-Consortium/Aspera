import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { MatTableDataSource, MatSort } from '@angular/material';
import { Transaction } from '../../../lib/model';
import { AccountService } from '../../../lib/services';

@Component({
    selector: 'app-transactions',
    styleUrls: ['./transactions.component.css'],
    templateUrl: './transactions.component.html'
})
export class TransactionsComponent {
    private dataSource: MatTableDataSource<Transaction>;
    private displayedColumns: string[];

    @ViewChild(MatSort) sort: MatSort;

    constructor(
        private accountService: AccountService
    ) {}

    public ngOnInit() {
        this.displayedColumns = ['id', 'sender', 'recipient', 'amount', 'fee', 'type', 'attachment', 'timestamp'];
        this.dataSource = new MatTableDataSource<Transaction>();

        this.accountService.getTransactions("14276304710264415646").then(transactions => {
            this.dataSource.data = transactions;
        })
    }

    public ngAfterViewInit() {
        this.dataSource.sort = this.sort;
    }

    public applyFilter(filterValue: string) {
        filterValue = filterValue.trim(); // Remove whitespace
        filterValue = filterValue.toLowerCase(); // MatTableDataSource defaults to lowercase matches
        this.dataSource.filter = filterValue;
    }
}
