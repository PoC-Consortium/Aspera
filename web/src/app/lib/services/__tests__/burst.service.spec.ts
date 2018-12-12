import {inject, TestBed} from '@angular/core/testing';
import {HttpClientTestingModule} from '@angular/common/http/testing';

import {BurstService} from '../burst.service';
import {I18nService} from '../../i18n/i18n.service';
import {StoreService} from '../store.service';
import {testConfigFactory} from "../../config/store.config";
import {BehaviorSubject} from "rxjs/BehaviorSubject";
import {Settings} from "../../model";

jest.mock('../../i18n/i18n.service');

describe('BurstService', () => {

    beforeEach(() => {

        const storeServiceMock = new StoreService(testConfigFactory());
        storeServiceMock.settings = new BehaviorSubject(new Settings());

        TestBed.configureTestingModule({
            imports: [HttpClientTestingModule],
            providers: [
                BurstService,
                I18nService,
                {provide: StoreService, useValue: storeServiceMock}
            ]
        });
    });

    it('should be created', inject([BurstService], (service: BurstService) => {
            expect(service).not.toBeNull();
        })
    );
});
