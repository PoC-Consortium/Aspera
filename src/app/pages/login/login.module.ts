import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';

import { routing }       from './login.routing';
import { SharedModule } from '../../lib/shared.module';
import { TranslateModule } from '@ngx-translate/core';

import { LoginComponent } from './login.component';
import { SetupModule } from '../dashboard/setup/setup.module';

@NgModule({
    imports: [
        CommonModule,
        SharedModule,
        TranslateModule,
        SetupModule,
        routing,
    ],
    declarations: [
        LoginComponent
    ]
})
export class LoginModule {
}
