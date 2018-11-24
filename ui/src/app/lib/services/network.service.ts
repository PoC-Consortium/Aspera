/*
* Copyright 2018 PoC-Consortium
*/

import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import 'rxjs/add/operator/timeout'

import { BurstNode, constants } from "../model";
import { NoConnectionError } from "../model/error";

/*
* NetworkService class
*
* Doing network stuff
*/
@Injectable()
export class NetworkService {

    constructor(
        private http: HttpClient
    ) {}

    public latency(node: BurstNode) : Promise<number> {
        return new Promise((resolve, reject) => {
            let timeStart: number = performance.now();
            return this.http.get(this.constructNodeUrl(node))
            .timeout(constants.connectionTimeout)
            .toPromise()
            .then(response => {
                let timeEnd: number = performance.now();
                resolve(timeEnd - timeStart);
            }).catch(e => {
                console.log(e)
                reject(new NoConnectionError("Connection timed out!"))
            });
        });
    }

    public constructNodeUrl(node: BurstNode) : string {
        return node.address + ":" + node.port + "/burst";
    }

}
