import { NgModule, ModuleWithProviders }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';

import { TimeAgoPipe } from 'time-ago-pipe';

import {
    AppFooterComponent,
    AppHeaderComponent,
    ControlSidebarComponent,
    MenuAsideComponent,
    NotificationBoxComponent
} from './components';

import {

} from './pipes';

import {
    AccountService,
    CryptoService,
    LoggerService,
    MarketService,
    NotificationService,
    StoreService
} from './services';

import {

} from './validators';

const NGA_COMPONENTS = [
    AppFooterComponent,
    AppHeaderComponent,
    ControlSidebarComponent,
    MenuAsideComponent,
    NotificationBoxComponent
];

const NGA_DIRECTIVES = [

];

const NGA_PIPES = [
    TimeAgoPipe
];

const NGA_SERVICES = [
    AccountService,
    CryptoService,
    LoggerService,
    MarketService,
    NotificationService,
    StoreService,
    TranslateService
];

const NGA_VALIDATORS = [

];

@NgModule({
    declarations: [
        ...NGA_PIPES,
        ...NGA_DIRECTIVES,
        ...NGA_COMPONENTS
    ],
    imports: [
        CommonModule,
        RouterModule,
        FormsModule,
        ReactiveFormsModule,
        TranslateModule.forRoot(),
        MatButtonModule,
        MatIconModule,
        MatMenuModule
    ],
    exports: [
        ...NGA_PIPES,
        ...NGA_DIRECTIVES,
        ...NGA_COMPONENTS
    ]
})
export class SharedModule {
    static forRoot(): ModuleWithProviders {
        return <ModuleWithProviders>{
            ngModule: SharedModule,
            providers: [
                ...NGA_VALIDATORS,
                ...NGA_SERVICES
            ],
        };
    }
}
