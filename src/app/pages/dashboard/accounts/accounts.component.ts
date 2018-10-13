import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { MatTableDataSource, MatSort } from '@angular/material';
import { Account } from '../../../lib/model';
import { AccountsListActions } from '../../../auth/actions';
import { AccountService, StoreService } from '../../../lib/services';
import { Store } from '@ngrx/store';

@Component({
    selector: 'app-accounts',
    styleUrls: ['./accounts.component.css'],
    templateUrl: './accounts.component.html'
})
export class AccountsComponent {
    private accounts: Account[];
    private dataSource: MatTableDataSource<Account>;
    private displayedColumns: string[];

    @ViewChild(MatSort) sort: MatSort;

    constructor(
        private storeService: StoreService,
        private accountService: AccountService
    ) {}

    public ngOnInit() {
        this.accounts = [];
        this.displayedColumns = ['id', 'address', 'alias', 'balance', 'selected'];
        this.dataSource = new MatTableDataSource<Account>();
  
        this.storeService.ready.subscribe((ready) => {
            this.storeService.getAllAccounts().then((accounts) => {
                this.accounts = accounts;
                this.dataSource.data = this.accounts;
            })
          });
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
