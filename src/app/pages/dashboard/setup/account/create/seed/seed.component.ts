import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';
import { CryptoService } from '../../../../../../lib/services';
import { CreateService } from '../create.service';

@Component({
    selector: 'app-account-create-seed',
    styleUrls: ['./seed.component.scss'],
    templateUrl: './seed.component.html'
})
export class AccountCreateSeedComponent implements OnInit {
    private seedLimit: number = 10;
    private seed: any[] = [];
    private update: boolean = false;
    private generating = false;
    private progress = 0;

    constructor(
        private router: Router,
        private cryptoService: CryptoService,
        private createService: CreateService
    ) { }

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

    public movement(e) {
        this.seed.push([e.clientX, e.clientY, new Date()]);
        if (!this.update) {
            this.update = true
            setTimeout(() => {
                this.progress = this.seed.length / this.seedLimit * 100;
                this.update = false;
            }, 100)
        }
        if (this.seed.length >= this.seedLimit) {
            this.generating = true;
            this.cryptoService.generatePassPhrase(this.seed)
                .then(phrase => {
                    this.createService.setPassphrase(phrase);
                    this.cryptoService.generateMasterKeys(this.createService.getCompletePassphrase()).then(keys => {
                        this.cryptoService.getAccountIdFromPublicKey(keys.publicKey).then(id => {
                            this.createService.setId(id);
                            this.cryptoService.getBurstAddressFromAccountId(id).then(address => {
                                this.createService.setAddress(address);
                                setTimeout(x => {
                                    this.createService.setStepIndex(1)
                                }, 1000);
                                this.generating = false;
                            })
                        })
                    })
                });
        }
    }

}
