import { Component, Input } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { MarketService } from '../../services';


@Component({
    selector: 'app-header',
    styleUrls: ['./app-header.component.scss'],
    templateUrl: './app-header.component.html'
})
export class AppHeaderComponent {

    constructor(
        private marketService: MarketService
    ) {
    }

}
