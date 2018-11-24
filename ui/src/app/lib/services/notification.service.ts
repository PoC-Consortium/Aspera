import { Injectable } from '@angular/core';

@Injectable()
export class NotificationService {
    constructor() { }

    public success = (body: string, title = 'Operation successful'): void => {

    }

    public error = (body: string, title = 'An error occured'): void => {

    }

    public warning = (body: string, title = 'Something went wrong'): void => {

    }
}
