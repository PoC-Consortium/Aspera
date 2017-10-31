import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SharedModule } from '../../lib/shared.module';

import { AccountsComponent } from './accounts.component';
import { routing } from './accounts.routing';

@NgModule({
    imports: [
        CommonModule,
        FormsModule,
        SharedModule,
        routing
    ],
    declarations: [
        AccountsComponent
    ],
    providers: [
    ]
})
export class AccountsModule { }
