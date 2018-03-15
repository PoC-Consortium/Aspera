import { Component, OnInit, OnDestroy, NgZone, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';
import { CryptoService } from '../../../../../../lib/services';
import { CreateService } from '../create.service';

@Component({
    selector: 'app-account-create-record',
    styleUrls: ['./record.component.scss'],
    templateUrl: './record.component.html'
})
export class AccountCreateRecordComponent implements OnInit {

    constructor(
        private router: Router,
        private createService: CreateService,
        private _ngZone: NgZone
    ) { }

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

    public reset() {
        this.createService.reset();
        // Angular Stepper hack
        setTimeout(x => {
            this.createService.setStepIndex(0)
        }, 0);
    }

    public next() {
        this.createService.setStepIndex(2);
    }

}
