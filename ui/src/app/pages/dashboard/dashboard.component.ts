import { Component, OnInit, ViewEncapsulation, Output, HostListener } from '@angular/core';
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
    @HostListener('window:resize', ['$event'])
    onResize(event) {
        console.log(event.target.innerWidth);
        this.opened = (event.target.innerWidth > 800);
    }
    constructor(
        private marketService: MarketService
    ) {

    }

    public ngOnInit() { 
        let timer = Observable.timer(2000, 60000);
        timer.subscribe(t =>
            this.marketService.updateCurrency()
        );

        
    } 


}
