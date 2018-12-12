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
    MenuAsideComponent,
    SendBurstDialogComponent
} from './components';

import {
    AccountService,
    CryptoService,
    LoggerService,
    MarketService,
    NotificationService,
    NetworkService,
    StoreService,
    BurstService,
} from './services';
import { I18nModule } from './i18n/i18n.module';
import { MatDialogModule, MatFormFieldModule, MatInputModule, MatTabsModule, MatCheckboxModule, MatGridListModule } from '@angular/material';
import { SendBurstFormComponent } from './components/send-burst-form';
import { SendMultiOutFormComponent } from './components/send-multi-out-form';
import { BurstInputValidatorDirective } from './components/send-burst-form/burst-input-validator.directive';

import { NgxMaskModule } from 'ngx-mask';
import { SendMessageFormComponent } from './components/send-message-form/send-message-form.component';

const NGA_COMPONENTS = [
    AppFooterComponent,
    AppHeaderComponent,
    MenuAsideComponent,
    SendBurstDialogComponent
];

const NGA_DIRECTIVES = [
    BurstInputValidatorDirective
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
    TranslateService,
    BurstService
];

const NGA_VALIDATORS = [

];

@NgModule({
    declarations: [
        ...NGA_PIPES,
        ...NGA_DIRECTIVES,
        ...NGA_COMPONENTS,
        SendBurstFormComponent,
        SendMultiOutFormComponent,
        SendMessageFormComponent
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
        MatFormFieldModule,
        MatInputModule,
        MatListModule,
        MatDialogModule,
        MatSidenavModule,
        MatCheckboxModule,
        MatGridListModule,
        I18nModule,
        MatTabsModule,
        NgxMaskModule.forRoot()
    ],
    entryComponents: [ SendBurstDialogComponent ],
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
