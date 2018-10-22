import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';
import { CryptoService } from '../../../../../../lib/services';
import { CreateService } from '../../create.service';
import { SetupService } from '../../../setup.service';

@Component({
    selector: 'app-account-create-pin',
    styleUrls: ['./pin.component.scss'],
    templateUrl: './pin.component.html'
})
export class AccountCreatePinComponent implements OnInit {

    constructor(
        private router: Router,
        private createService: CreateService,
        private setupService: SetupService
    ) { }

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

    public back() {
        this.createService.setStepIndex(1);
    }

    public finish(pin: string) {
        this.createService.setPin(pin);
        this.createService.createActiveAccount();
    }

}
