import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

@Injectable()
export class CreateService {
    private passphrase: string[];
    private id: string;
    private address: string;
    private stepIndex: number;

    constructor() {
        this.stepIndex = 0;
        this.reset();
    }

    public setPassphrase(passphrase: string[]) {
        this.passphrase = passphrase;
    }

    public getPassphrase(): string[] {
        return this.passphrase;
    }

    public getPassphrasePart(index: number): string {
        return this.passphrase[index];
    }

    public getCompletePassphrase(): string {
        return this.passphrase.join(" ");
    }

    public setId(id: string) {
        this.id = id;
    }

    public getId(id: string) {
        return this.id;
    }

    public setAddress(address: string) {
        this.address = address;
    }

    public getAddress() : string {
        return this.address;
    }

    public setStepIndex(index: number) {
        this.stepIndex = index;
    }

    public getStepIndex() : number {
        return this.stepIndex;
    }

    public isPassphraseGenerated() : boolean {
        return this.passphrase.length > 0 && this.address != undefined && this.id != undefined
    }

    public reset() {
        this.passphrase = [];
        this.id = undefined;
        this.address = undefined;
    }
}
