import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';

import { routing }       from './login.routing';
import { SharedModule } from '../../lib/shared.module';
import { TranslateModule } from '@ngx-translate/core';

import { LoginComponent } from './login.component';
import { SetupModule } from '../dashboard/setup/setup.module';
import { MatIconModule, MatInputModule } from '@angular/material';
import { NotifierModule } from 'angular-notifier';
import { FormsModule } from '@angular/forms';
import { I18nModule } from '../../lib/i18n/i18n.module';
import { NgxMaskModule } from 'ngx-mask';

@NgModule({
    imports: [
        CommonModule,
        SharedModule,
        TranslateModule,
        SetupModule,
        routing,
        MatIconModule,
        NotifierModule,
        FormsModule,
        I18nModule,
        MatInputModule,
        NgxMaskModule.forRoot()
    ],
    declarations: [
        LoginComponent
    ]
})
export class LoginModule {
}
