import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { MatPaginator, MatTableDataSource, MatSort } from '@angular/material';
import { TimeAgoPipe } from 'time-ago-pipe';
import { Transaction } from '../../../lib/model';
import { AccountService } from '../../../lib/services';
import { Converter } from '../../../lib/util';

@Component({
    selector: 'app-transactions',
    styleUrls: ['./transactions.component.css'],
    templateUrl: './transactions.component.html'
})
export class TransactionsComponent {
    private dataSource: MatTableDataSource<Transaction>;
    private displayedColumns: string[];
    private account = 'BURST-RFR4-DQ4V-XF3Q-9FD3E';

    private incoming: boolean;
    private outgoing: boolean;

    @ViewChild(MatPaginator) paginator: MatPaginator;
    @ViewChild(MatSort) sort: MatSort;

    constructor(
        private accountService: AccountService
    ) {}

    public ngOnInit() {
        this.displayedColumns = ['type', 'opposite', 'amount', 'fee', 'attachment', 'timestamp', 'confirmed'];
        this.dataSource = new MatTableDataSource<Transaction>();

        this.accountService.getTransactions("8504323719802107618").then(transactions => {
            this.dataSource.data = transactions;
        })
    }

    public ngAfterViewInit() {
        this.dataSource.paginator = this.paginator;
        this.dataSource.sort = this.sort;
    }

    public applyFilter(filterValue: string) {
        filterValue = filterValue.trim(); // Remove whitespace
        filterValue = filterValue.toLowerCase(); // MatTableDataSource defaults to lowercase matches
        this.dataSource.filter = filterValue;
    }

    public isOwnAccount(address: string): boolean {
        return address != undefined && address == this.account;
    }

    public convertTimestamp(timestamp: number): Date {
        return Converter.convertTimestampToDate(timestamp);
    }
}
