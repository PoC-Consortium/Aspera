import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

import { Account, Node } from '../../../lib/model';

@Injectable()
export class SetupService {
    private account: Account;
    private node: Node;
    private stepIndex: number;

    constructor() {
        this.reset();
    }

    public setNode(node: Node) {
        this.node = node;
    }

    public getNode(): Node {
        return this.node;
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
        this.node = undefined;
        this.stepIndex = 0;
    }
}
