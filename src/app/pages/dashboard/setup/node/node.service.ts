import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

@Injectable()
export class NodeService {
    private address: string;
    private port: number;

    constructor() {
        this.reset();
    }

    public setAddress(address: string) {
        this.address = address;
    }

    public getAddress() : string {
        return this.address;
    }

    public setPort(port: number) {
        this.port = port;
    }

    public getPort() : number {
        return this.port;
    }

    public reset() {
        this.address = undefined;
        this.port = 0;
    }
}
