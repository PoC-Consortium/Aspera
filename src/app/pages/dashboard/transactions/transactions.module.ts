import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SharedModule } from '../../../lib/shared.module';

import { TransactionsComponent } from './transactions.component';
import { TransactionsRouting } from './transactions.routing';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule, MatSortModule } from '@angular/material';
import { MatTableModule } from '@angular/material/table';

@NgModule({
    imports: [
        TransactionsRouting,
        CommonModule,
        FormsModule,
        SharedModule,
        MatCardModule,
        MatFormFieldModule,
        MatGridListModule,
        MatIconModule,
        MatInputModule,
        MatSortModule,
        MatTableModule
    ],
    declarations: [
        TransactionsComponent
    ],
    providers: [
    ]
})
export class TransactionsModule { }
