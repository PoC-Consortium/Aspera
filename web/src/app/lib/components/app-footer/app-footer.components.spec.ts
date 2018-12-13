import {
    async,
    ComponentFixture,
    TestBed
} from '@angular/core/testing';

import { AppFooterComponent } from './app-footer.component';

describe('AppFooterComponent', () => {
    let component: AppFooterComponent;
    let fixture: ComponentFixture<AppFooterComponent>;
    let dom : any;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [ AppFooterComponent ]
        });
        fixture = TestBed.createComponent(AppFooterComponent);
        component = fixture.componentInstance;
    }));

    test('should exist', () => {
        expect(component).toBeDefined();
    });

    test('should be rendered correctly', () => {

        expect(fixture).toMatchSnapshot();
    });
});
