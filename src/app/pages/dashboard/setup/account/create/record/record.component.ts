import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
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
        private createService: CreateService
    ) { }

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

}
