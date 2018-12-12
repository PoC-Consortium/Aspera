import {async, inject, TestBed} from '@angular/core/testing';

import {StoreService} from '../store.service';
import {StoreConfig, testConfigFactory} from "../../config/store.config";
import {Settings} from "../../model";

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

            const defaultSettings = new Settings();

            expect(settings).not.toBeNull();
            expect(settings).toEqual({...defaultSettings} );

        })
        )
    );

});
