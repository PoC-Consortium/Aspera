import {async, inject, TestBed} from '@angular/core/testing';

import {StoreService} from '../store.service';
import {StoreConfig, testConfigFactory} from "../../config/store.config";

describe('StoreService', () => {

    beforeEach(() => {

        TestBed.configureTestingModule({
            providers: [
                StoreService,
                {provide: StoreConfig, useFactory: testConfigFactory}
            ]
        });
    });

    it('should be created', inject([StoreService], (service: StoreService) => {
            expect(service).toBeTruthy();
        })
    );

    it('should be initialized with default values', async(inject([StoreService], async (service: StoreService) => {
            service.init();
            const settings = await service.getSettings();

            const expectedSettings = {
                contacts: [],
                currency: "USD",
                id: "settings",
                language: "en",
                marketUrl: "http://localhost:4200/v1/ticker/burst/",
                node: "http://localhost:4200/burst",
                theme: "light",
                version: "0.2.1",
            };
            expect(settings).not.toBeNull();
            expect(settings).toEqual(expectedSettings);

        })
        )
    );

});
