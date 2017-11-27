import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SharedModule } from '../../../lib/shared.module';

import { HomeComponent } from './home.component';
import { routing } from './home.routing';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule, MatSortModule } from '@angular/material';
import { MatTableModule } from '@angular/material/table';

@NgModule({
    imports: [
        CommonModule,
        FormsModule,
        SharedModule,
        routing,
        MatCardModule,
        MatFormFieldModule,
        MatInputModule,
        MatSortModule,
        MatTableModule
    ],
    declarations: [
        HomeComponent
    ],
    providers: [
    ]
})
export class HomeModule { }
