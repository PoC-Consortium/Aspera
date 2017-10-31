import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SharedModule } from '../../lib/shared.module';

import { CreateComponent } from './create.component';

import { NotificationService } from '../../lib/services';
import { routing } from './create.routing';

@NgModule({
    imports: [
        CommonModule,
        FormsModule,
        SharedModule,
        routing
    ],
    declarations: [
        CreateComponent
    ],
    providers: [
        NotificationService
    ]
})
export class CreateModule { }
