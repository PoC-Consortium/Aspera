import { Component, OnInit } from '@angular/core';
import { StoreService, AccountService, BurstService } from './lib/services';
import { GetBlockchainStatusResponse, Account } from './lib/model';
import { map } from 'rxjs-compat/operator/map';

@Component( {
    styleUrls: ['app.component.scss'],
    selector: 'app',
    templateUrl: './app.component.html'
})
export class App {
    constructor(private storeService: StoreService,
                private accountService: AccountService,
                private burstService: BurstService) {
    }

    firstTime = true;
    isScanning = false;
    downloadingBlockchain = false;
    previousLastBlock = "0";
    lastBlock = "0";
    account: Account;
    accounts: Account[];

    ngOnInit() {
        this.storeService.ready.subscribe((ready) => {
            if (ready) {
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
                            this.accountService.synchronizeAccount(account);
                        }, 1);
                    })
                })
                // todo
                // this.burstService.getBlockchainStatus().then((response: GetBlockchainStatusResponse) => {
                //     this.previousLastBlock = this.firstTime && response.lastBlock;

                //     if (this.firstTime) {
                //         this.firstTime = false;
                //         return this.burstService.getBlock(response.lastBlock).then(() => {
                //             // handleInitialBlocks
                //         });
                //     } else if (response.isScanning) { // todo: confusing
                //         this.isScanning = true;
                //     } else if (this.isScanning) {
                //         this.isScanning = false;
                //         this.burstService.resetBlocks();
                //         this.burstService.resetTempBlocks();
                //         this.burstService.getBlock(response.lastBlock).then(() => {
                //             //handleInitialBlocks
                //         });
                //         if (this.account) {
                //             this.accountService.synchronizeAccount(this.account);
                //         }
                //     } else if (this.previousLastBlock !== response.lastBlock) {
                //         this.burstService.resetTempBlocks();
                //         this.burstService.getBlock(response.lastBlock).then(() => {
                //             //handleNewBlocks
                //         });
                //         if (this.account) {
                //             this.accountService.synchronizeAccount(this.account);
                //         }
                //     } else {
                //         if (this.account) {
                //             this.accountService.getUnconfirmedTransactions(this.account.id).then((unconfirmedTransactions) => {
                //                 this.accountService.handleIncomingTransactions(unconfirmedTransactions)
                //             })
                //         }
                //         if (this.downloadingBlockchain) {
                //             this.updateBlockchainDownloadProgress();
                //         }
                //     }
                // });
            }
        });

    }

    updateBlockchainDownloadProgress() {}

}
