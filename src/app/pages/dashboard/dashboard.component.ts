import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { LoggerService } from '../../lib/services';
import { Observable, Subject } from 'rxjs/Rx';
import { MarketService } from '../../lib/services';

@Component({
    selector: 'particle-dashboard',
    styles: [],
    templateUrl: './dashboard.html'
})
export class DashboardComponent implements OnInit {
    private mylinks: Array<any> = [];

    constructor(
        private marketService: MarketService
    ) {

    }

    public ngOnInit() {
        let ie = this.detectIE();
        if (!ie) {
            window.dispatchEvent(new Event('resize'));
        } else {
            // solution for IE from @hakonamatata
            let event = document.createEvent('Event');
            event.initEvent('resize', false, true);
            window.dispatchEvent(event);
        }

        // define here your own links menu structure
        this.mylinks = [
            {
                'title': 'Accounts',
                'icon': 'account_box',
                'link': ['/accounts']
            },
            {
                'title': 'Transactions',
                'icon': 'account_balance_wallet',
                'link': ['/transactions']
            },
            {
                'title': 'Assets',
                'icon': 'clear_all',
                'link': ['/transactions']
            }
        ];

        let timer = Observable.timer(2000, 10000);
        timer.subscribe(t =>
            this.marketService.updateCurrency()
        );
    }

    protected detectIE(): any {
        let ua = window.navigator.userAgent;

        // Test values; Uncomment to check result …
        // IE 10
        // ua = 'Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Trident/6.0)';
        // IE 11
        // ua = 'Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko';
        // IE 12 / Spartan
        // ua = 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36 Edge/12.0';
        // Edge (IE 12+)
        // ua = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)
        // Chrome/46.0.2486.0 Safari/537.36 Edge/13.10586';

        let msie = ua.indexOf('MSIE ');
        if (msie > 0) {
            // IE 10 or older => return version number
            return parseInt(ua.substring(msie + 5, ua.indexOf('.', msie)), 10);
        }

        let trident = ua.indexOf('Trident/');
        if (trident > 0) {
            // IE 11 => return version number
            let rv = ua.indexOf('rv:');
            return parseInt(ua.substring(rv + 3, ua.indexOf('.', rv)), 10);
        }

        let edge = ua.indexOf('Edge/');
        if (edge > 0) {
            // Edge (IE 12+) => return version number
            return parseInt(ua.substring(edge + 5, ua.indexOf('.', edge)), 10);
        }

        // other browser
        return false;
    }


}
