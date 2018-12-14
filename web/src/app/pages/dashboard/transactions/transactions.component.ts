import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { MatPaginator, MatTableDataSource, MatSort } from '@angular/material';
import { TimeAgoPipe } from 'time-ago-pipe';
import { Transaction, Account } from '../../../lib/model';
import { AccountService, StoreService } from '../../../lib/services';
import { Converter } from '../../../lib/util';

@Component({
    selector: 'app-transactions',
    styleUrls: ['./transactions.component.css'],
    templateUrl: './transactions.component.html'
})
export class TransactionsComponent {
    private dataSource: MatTableDataSource<Transaction>;
    private displayedColumns: string[];

    private incoming: boolean = true;
    private outgoing: boolean = true;
    private account: Account;

    @ViewChild(MatPaginator) paginator: MatPaginator;
    @ViewChild(MatSort) sort: MatSort;

    constructor(
        private accountService: AccountService,
        private storeService: StoreService
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
        this.dataSource.paginator = this.paginator;
        this.dataSource.sort = this.sort;
    }

    public applyFilter(filterValue: string) {
        filterValue = filterValue.trim(); // Remove whitespace
        filterValue = filterValue.toLowerCase(); // MatTableDataSource defaults to lowercase matches
        this.dataSource.filter = filterValue;
    }

    public isOwnAccount(address: string): boolean {
        return address != undefined && address == this.account.address;
    }

    public convertTimestamp(timestamp: number): Date {
        return Converter.convertTimestampToDate(timestamp);
    }

    public getTransactionNameFromType(transaction: Transaction) {
        var transactionType = "unknown";
        if (transaction.type === 0) {
            transactionType = "ordinary_payment";
        }
        else if (transaction.type == 1) {
            switch (transaction.subtype) {
            case 0:
                transactionType = "arbitrary_message";
                break;
            case 1:
                transactionType = "alias_assignment";
                break;
            case 2:
                transactionType = "poll_creation";
                break;
            case 3:
                transactionType = "vote_casting";
                break;
            case 4:
                transactionType = "hub_announcements";
                break;
            case 5:
                transactionType = "account_info";
                break;
            case 6:
                if (transaction.attachment.priceNQT == "0") {
                    if (transaction.senderId == this.account.address && transaction.recipientId == this.account.address) {
                        transactionType = "alias_sale_cancellation";
                    }
                    else {
                        transactionType = "alias_transfer";
                    }
                }
                else {
                    transactionType = "alias_sale";
                }
                break;
            case 7:
                transactionType = "alias_buy";
                break;
            }
        }
        else if (transaction.type == 2) {
            switch (transaction.subtype) {
            case 0:
                transactionType = "asset_issuance";
                break;
            case 1:
                transactionType = "asset_transfer";
                break;
            case 2:
                transactionType = "ask_order_placement";
                break;
            case 3:
                transactionType = "bid_order_placement";
                break;
            case 4:
                transactionType = "ask_order_cancellation";
                break;
            case 5:
                transactionType = "bid_order_cancellation";
                break;
            }
        }
        else if (transaction.type == 3) {
            switch (transaction.subtype) {
            case 0:
                transactionType = "marketplace_listing";
                break;
            case 1:
                transactionType = "marketplace_removal";
                break;
            case 2:
                transactionType = "marketplace_price_change";
                break;
            case 3:
                transactionType = "marketplace_quantity_change";
                break;
            case 4:
                transactionType = "marketplace_purchase";
                break;
            case 5:
                transactionType = "marketplace_delivery";
                break;
            case 6:
                transactionType = "marketplace_feedback";
                break;
            case 7:
                transactionType = "marketplace_refund";
                break;
            }
        }
        else if (transaction.type == 4) {
            switch (transaction.subtype) {
            case 0:
                transactionType = "balance_leasing";
                break;
            }
        }
        else if (transaction.type == 20) {
            switch (transaction.subtype) {
            case 0:
                transactionType = "Reward Recipient Assignment";
                break;
            }
        }
        else if (transaction.type == 21) {
            switch (transaction.subtype) {
            case 0:
                transactionType = "Escrow Creation";
                break;
            case 1:
                transactionType = "Escrow Signing";
                break;
            case 2:
                transactionType = "Escrow Result";
                break;
            case 3:
                transactionType = "Subscription Subscribe";
                break;
            case 4:
                transactionType = "Subscription Cancel";
                break;
            case 5:
                transactionType = "Subscription Payment";
                break;
            }
        }
        else if (transaction.type == 22) {
            switch (transaction.subtype) {
            case 0:
                transactionType = "AT Creation";
                break;
            case 1:
                transactionType = "AT Payment";
                break;
            }
        }
        return transactionType;
    };
}
