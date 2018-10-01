import { Component, OnInit, ViewEncapsulation, Output } from '@angular/core';
import { LoggerService, StoreService, AccountService } from '../../lib/services';
import { Observable, Subject } from 'rxjs/Rx';
import { MarketService } from '../../lib/services';
import { Transaction } from '../../lib/model';

@Component({
    selector: 'particle-dashboard',
    styleUrls: ['./dashboard.component.scss'],
    templateUrl: './dashboard.component.html'
})
export class DashboardComponent implements OnInit {
    opened: boolean = true;
    events: string[] = [];
    @Output() recentTransactions: Transaction[];

    constructor(
        private marketService: MarketService
    ) {

    }

    public ngOnInit() {
        // let timer = Observable.timer(2000, 10000);
        // timer.subscribe(t =>
        //     this.marketService.updateCurrency()
        // );
    } 


}
