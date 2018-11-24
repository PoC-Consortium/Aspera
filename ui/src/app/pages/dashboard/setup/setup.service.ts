import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

import { Account, BurstNode } from '../../../lib/model';

@Injectable()
export class SetupService {
    private account: Account;
    private BurstNode: BurstNode;
    private stepIndex: number;

    constructor() {
        this.reset();
    }

    public setBurstNode(BurstNode: BurstNode) {
        this.BurstNode = BurstNode;
    }

    public getBurstNode(): BurstNode {
        return this.BurstNode;
    }

    public setAccount(account: Account) {
        this.account = account;
    }

    public getAccount() : Account {
        return this.account;
    }

    public setStepIndex(index: number) {
        this.stepIndex = index;
    }

    public getStepIndex() : number {
        return this.stepIndex;
    }

    public reset() {
        this.account = undefined;
        this.BurstNode = undefined;
        this.stepIndex = 0;
    }
}
