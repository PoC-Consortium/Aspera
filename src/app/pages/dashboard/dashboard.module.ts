import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';

import { routing }       from './dashboard.routing';
import { SharedModule } from '../../lib/shared.module';
import { TranslateModule } from '@ngx-translate/core';

import { DashboardComponent } from './dashboard.component';
import { MatSidenavModule } from '@angular/material';

@NgModule({
    imports: [
        CommonModule,
        SharedModule,
        TranslateModule,
        MatSidenavModule,
        routing
    ],
    declarations: [
        DashboardComponent
    ]
})
export class DashboardModule {
}
