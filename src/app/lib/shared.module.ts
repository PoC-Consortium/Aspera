import { NgModule, ModuleWithProviders }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';

import { TimeAgoPipe } from 'time-ago-pipe';

import {
    AppFooterComponent,
    AppHeaderComponent,
    ControlSidebarComponent,
    MenuAsideComponent,
    NotificationBoxComponent
} from './components';

import {
    AccountService,
    CryptoService,
    LoggerService,
    MarketService,
    NotificationService,
    NetworkService,
    StoreService,
} from './services';
import { I18nModule } from './i18n/i18n.module';

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
    NetworkService,
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
        MatMenuModule,
        MatListModule,
        MatSidenavModule,
        I18nModule
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
