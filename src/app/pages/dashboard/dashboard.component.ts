import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { LoggerService } from '../../lib/services';
import { Observable, Subject } from 'rxjs/Rx';
import { MarketService } from '../../lib/services';

@Component({
    selector: 'particle-dashboard',
    styleUrls: ['./dashboard.component.scss'],
    templateUrl: './dashboard.component.html'
})
export class DashboardComponent implements OnInit {
    opened: boolean = true;
    events: string[] = [];

    constructor(
        private marketService: MarketService
    ) {
    }

    public ngOnInit() {
        let timer = Observable.timer(2000, 10000);
        timer.subscribe(t =>
            this.marketService.updateCurrency()
        );
    } 


}
