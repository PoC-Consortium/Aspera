import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SharedModule } from '../../../lib/shared.module';

import { HomeComponent } from './home.component';
import { routing } from './home.routing';

@NgModule({
    imports: [
        CommonModule,
        FormsModule,
        SharedModule,
        routing
    ],
    declarations: [
        HomeComponent
    ],
    providers: [
    ]
})
export class HomeModule { }
