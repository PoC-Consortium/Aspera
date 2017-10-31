import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';

import { routing }       from './dashboard.routing';
import { SharedModule } from '../../lib/shared.module';
import { TranslateModule } from '@ngx-translate/core';

import { DashboardComponent } from './dashboard.component';

@NgModule({
    imports: [
        CommonModule,
        SharedModule,
        TranslateModule,
        routing
    ],
    declarations: [
        DashboardComponent
    ]
})
export class DashboardModule {
}
