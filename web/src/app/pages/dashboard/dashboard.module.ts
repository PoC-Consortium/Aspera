import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';

import { routing }       from './dashboard.routing';
import { SharedModule } from '../../lib/shared.module';
import { TranslateModule } from '@ngx-translate/core';

import { DashboardComponent } from './dashboard.component';
import { MatSidenavModule } from '@angular/material';
import { AccountsModule } from './accounts/accounts.module';
import { SetupModule } from './setup/setup.module';
import { HomeModule } from './home/home.module';
import { TransactionsModule } from './transactions/transactions.module';
import { MessagesComponent } from './messages/messages.component';

@NgModule({
    imports: [
        CommonModule,
        SharedModule,
        TranslateModule,
        MatSidenavModule,
        AccountsModule,
        SetupModule,
        HomeModule,
        TransactionsModule,
        routing,
    ],
    declarations: [
        DashboardComponent,
        MessagesComponent
    ]
})
export class DashboardModule {
}
