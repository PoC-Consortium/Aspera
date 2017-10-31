import { Injectable, Input } from '@angular/core';
import { environment } from '../../../environments/environment';

@Injectable()
export class LoggerService {

    constructor() {
    }

    public log(component: string, msg?: string) {
        if (!environment.silent) {
            console.log(component + ': ' + msg);
        }
    }
}
