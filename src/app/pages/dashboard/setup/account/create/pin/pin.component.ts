import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';
import { CryptoService } from '../../../../../../lib/services';
import { CreateService } from '../create.service';

@Component({
    selector: 'app-account-create-pin',
    styleUrls: ['./pin.component.scss'],
    templateUrl: './pin.component.html'
})
export class AccountCreatePinComponent implements OnInit {

    constructor(
        private router: Router,
        private createService: CreateService
    ) { }

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

    public reset() {
        this.createService.reset();
    }

}
