import { Component, OnInit } from '@angular/core';
import { StoreService, AccountService, BurstService } from './lib/services';
import { GetBlockchainStatusResponse, Account, HttpError, Block } from './lib/model';
import { map } from 'rxjs-compat/operator/map';
import { NotifierService } from 'angular-notifier';

@Component( {
    styleUrls: ['app.component.scss'],
    selector: 'app',
    templateUrl: './app.component.html'
})
export class App {
    constructor(private storeService: StoreService,
                private accountService: AccountService,
                private burstService: BurstService,
                private notificationService: NotifierService) {
    }

    firstTime = true;
    isScanning = false;
    downloadingBlockchain = false;
    previousLastBlock = "0";
    lastBlock = "0";
    account: Account;
    accounts: Account[];
    BLOCKCHAIN_STATUS_INTERVAL = 10000;

    ngOnInit() {
        this.storeService.ready.subscribe((ready) => {
            if (ready) {
                this.checkBlockchainStatus();
                setInterval(this.checkBlockchainStatus(), this.BLOCKCHAIN_STATUS_INTERVAL);
            }
        });

    }

    private checkBlockchainStatus(): (...args: any[]) => void {
        return () => {
            this.burstService.getBlockchainStatus().subscribe((response: GetBlockchainStatusResponse | HttpError | any) => {
                this.isScanning = !this.firstTime && (this.previousLastBlock != response.lastBlock);
                this.previousLastBlock = response.lastBlock;
                this.firstTime = false;
                if (this.isScanning && !this.firstTime) {
                    this.updateAccounts();
                    this.burstService.getBlock(response.lastBlock).subscribe((blockResponse: Block | HttpError | any) => {
                        if (blockResponse.errorCode) {
                            return this.notificationService.notify('error', this.burstService.translateServerError(response));
                        }
                    });
                } else if (this.account) {
                    this.accountService.synchronizeAccount(this.account).catch((error) => {
                        this.notificationService.notify('error', error.toString());
                    });
                }
            });
        };
    }

    private updateAccounts() {
        this.storeService.getSelectedAccount().then((account) => {
            if (account) {
                this.account = account;
                this.accountService.selectAccount(account);
            }
        });
        this.storeService.getAllAccounts().then((accounts) => {
            this.accounts = accounts;
            accounts.map((account) => {
                setTimeout(() => {
                    this.accountService.synchronizeAccount(account).catch((error) => {
                        this.notificationService.notify('error', error.toString());
                    })
                }, 1);
            });
        });
    }
}
