import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SharedModule } from '../../../lib/shared.module';

import { HomeComponent } from './home.component';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatInputModule, MatSortModule, MatIconModule } from '@angular/material';
import { MatTableModule } from '@angular/material/table';

@NgModule({
    imports: [
        CommonModule,
        FormsModule,
        SharedModule,
        MatCardModule,
        MatFormFieldModule,
        MatGridListModule,
        MatInputModule,
        MatIconModule,
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
