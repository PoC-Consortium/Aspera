import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { MatPaginator, MatTableDataSource, MatSort, MatDialog } from '@angular/material';
import { TimeAgoPipe } from 'time-ago-pipe';
import { Transaction, Account } from '../../../lib/model';
import { AccountService, StoreService, BurstService } from '../../../lib/services';
import { Converter } from '../../../lib/util';
import { TransactionDetailsDialogComponent } from './transaction-details-dialog';
import { FormControl } from '@angular/forms';

@Component({
    selector: 'app-transactions',
    styleUrls: ['./transactions.component.css'],
    templateUrl: './transactions.component.html'
})
export class TransactionsComponent {
    private dataSource: MatTableDataSource<Transaction>;
    private displayedColumns: string[];
    private account: Account;
    pickerFromField = new FormControl();
    pickerToField = new FormControl();


    @ViewChild(MatPaginator) paginator: MatPaginator;
    @ViewChild(MatSort) sort: MatSort;

    constructor(
        private accountService: AccountService,
        private storeService: StoreService,
        public burstService: BurstService,
        private dialog: MatDialog
    ) {}

    public ngOnInit() {
        this.displayedColumns = ['transaction_id', 'attachment', 'timestamp', 'type', 'amount', 'fee', 'account', 'confirmed'];
        this.dataSource = new MatTableDataSource<Transaction>();
        this.storeService.getSelectedAccount()
            .then((account) => {
                this.account = account;
                this.accountService.getTransactions(account.id).then(transactions => {
                    this.dataSource.data = transactions;
                })
            })
    } 

    public ngAfterViewInit() {
        const defaultFilterPredicate = this.dataSource.filterPredicate;
        this.dataSource.paginator = this.paginator;
        this.dataSource.sort = this.sort;
        this.dataSource.filterPredicate = (data, filterValue: string) => {
            let withinRange = true;
            if (this.pickerFromField.value && this.pickerToField.value) {
                withinRange = this.convertTimestamp(data.timestamp) > this.pickerFromField.value &&
                this.convertTimestamp(data.timestamp) < this.pickerToField.value;
            } else if (this.pickerFromField.value) {
                withinRange = this.convertTimestamp(data.timestamp) > this.pickerFromField.value;
            } else if (this.pickerToField.value) {
                withinRange = this.convertTimestamp(data.timestamp) < this.pickerToField.value;
            }
            return withinRange && defaultFilterPredicate(data, filterValue);
        }
    }

    public applyFilter(filterValue: string) {
        filterValue = filterValue.trim(); // Remove whitespace
        filterValue = filterValue.toLowerCase(); // MatTableDataSource defaults to lowercase matches
        this.dataSource.filter = filterValue || "burst";
    }

    public isOwnAccount(address: string): boolean {
        return address != undefined && address == this.account.address;
    }

    public convertTimestamp(timestamp: number): Date {
        return Converter.convertTimestampToDate(timestamp);
    }

    public openTransactionModal(transaction: Transaction): void {
        console.log(transaction);

        const dialogRef = this.dialog.open(TransactionDetailsDialogComponent, {
            width: '600px',
            data: { transaction: transaction }
        });

        dialogRef.afterClosed().subscribe(result => {
            console.log('The dialog was closed', result);
        });

    }
}
