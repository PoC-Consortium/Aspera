import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { MatTableDataSource, MatSort } from '@angular/material';
import { Account } from '../../../lib/model';

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

    public ngOnInit() {
        this.accounts = [];
        this.displayedColumns = ['id', 'alias', 'address', 'balance', 'selected'];
        this.dataSource = new MatTableDataSource<Account>();

        let a = new Account();
        a.id = "3664284351270490431";
        a.address = "BURST-26BZ-8J6S-7MFT-59KQ7";
        a.alias = "alias 1";
        a.balance = 624.36;
        a.selected = true;
        this.accounts.push(a);

        let b = new Account();
        b.id = "10416422395811764075";
        b.address = "BURST-K7VD-ZDYE-NAAJ-BJN53";
        b.alias = "alias 2";
        b.balance = 23000;
        b.selected = false;
        this.accounts.push(b);

        let c = new Account();
        c.id = "11642189530862802431";
        c.address = "BURST-LPHZ-QA3F-4RDF-CXC65";
        c.alias = "alias 3";
        c.balance = 2621.11205214;
        c.selected = false;
        this.accounts.push(c);

        this.dataSource.data = this.accounts;
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
