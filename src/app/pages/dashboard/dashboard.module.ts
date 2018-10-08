import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';

import { routing }       from './dashboard.routing';
import { SharedModule } from '../../lib/shared.module';
import { TranslateModule } from '@ngx-translate/core';

import { DashboardComponent } from './dashboard.component';
import { MatSidenavModule } from '@angular/material';
import { AccountsModule } from './accounts/accounts.module';
import { SetupModule } from './setup/setup.module';

@NgModule({
    imports: [
        CommonModule,
        SharedModule,
        TranslateModule,
        MatSidenavModule,
        AccountsModule,
        SetupModule,
        routing
    ],
    declarations: [
        DashboardComponent
    ]
})
export class DashboardModule {
}
