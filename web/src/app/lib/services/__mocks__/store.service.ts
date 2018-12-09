import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

import { Account, Settings} from "../../model";


// TODO: not sure if we need this mock...should be possible to do this another way (ngMocks, jest.mocking?)
@Injectable()
export class StoreService {

    private store: any;
    public ready: BehaviorSubject<any> = new BehaviorSubject(false);
    public settings: BehaviorSubject<any> = new BehaviorSubject(false);

    constructor() {}

    /*
    * Called on db start
    */
    public init() {}

    public setReady(state: boolean) {
        this.ready.next(state);
    }

    public setSettings(state: Settings) {
        this.settings.next(state);
    }

    /*
    * Method reponsible for saving/updating Account objects to the database.
    */
    public saveAccount(account: Account): Promise<Account> {
        return Promise.resolve(account)
    }

    /*
    * Method reponsible for getting the selected account from the database.
    */
    public getSelectedAccount(): Promise<Account> {
        return Promise.resolve(null)
    }

    /*
    * Method reponsible for selecting a new Account.
    */
    public selectAccount(account: Account): Promise<Account> {
        return Promise.resolve(account);
    }

    /*
    * Method reponsible for fetching all accounts from the database.
    */
    public getAllAccounts(): Promise<Account[]> {
        return Promise.resolve([])
    }

    /*
    * Method reponsible for finding an account by its numeric id from the database.
    */
    public findAccount(id: string): Promise<Account> {
        return Promise.resolve(new Account() );
    }

    /*
    * Method reponsible for removing an account from the database.
    */
    public removeAccount(account: Account): Promise<boolean> {
        return Promise.resolve(true);
    }

    /*
    * Method reponsible for saving/updating the global Settings object to the database.
    */
    public saveSettings(save: Settings): Promise<Settings> {
        return Promise.resolve(save);
    }

    /*
    * Method reponsible for fetching the global Settings object from the database.
    */
    public getSettings(): Promise<Settings> {
        return Promise.resolve(new Settings());
    }
}
