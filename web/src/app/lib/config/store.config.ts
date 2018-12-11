import {constants} from "../model";
import {LokiLocalStorageAdapter} from "lokijs" //loki-indexed-adapter.js";

export class StoreConfig {
    public databaseName: string;
    public persistenceAdapter: any;
}

const appConfigFactory = () => {
    const config = new StoreConfig();
    config.databaseName = constants.database;
    config.persistenceAdapter = new LokiLocalStorageAdapter();
    return config;
};

const testConfigFactory = () => {
    const config = new StoreConfig();
    config.databaseName = 'loki.test.db';
    config.persistenceAdapter = null; // in memory database
    return config;
};

export {
    testConfigFactory,
    appConfigFactory
}
