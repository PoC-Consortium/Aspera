import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { SharedModule } from '../../../lib/shared.module';

import { TransactionsComponent } from './transactions.component';

import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatDividerModule } from '@angular/material/divider';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon';
import { MatNativeDateModule, MatInputModule, MatSortModule, MatDialogModule } from '@angular/material';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatTabsModule } from '@angular/material/tabs';
import { I18nModule } from 'src/app/lib/i18n/i18n.module';
import { TransactionDetailsDialogComponent } from './transaction-details-dialog/transaction-details-dialog.component';
import { TransactionRowValueCellComponent } from './transaction-row-value-cell/transaction-row-value-cell.component';

@NgModule({
    imports: [
        CommonModule,
        FormsModule,
        SharedModule,
        MatButtonModule,
        MatCardModule,
        MatCheckboxModule,
        MatChipsModule,
        MatDatepickerModule,
        MatDividerModule,
        MatExpansionModule,
        MatFormFieldModule,
        MatGridListModule,
        MatIconModule,
        MatInputModule,
        MatNativeDateModule,
        MatPaginatorModule,
        MatSelectModule,
        MatSortModule,
        MatTableModule,
        MatTabsModule,
        MatTableModule,
        ReactiveFormsModule,
        I18nModule,
        MatDialogModule
    ],
    declarations: [
        TransactionDetailsDialogComponent,
        TransactionsComponent,
        TransactionRowValueCellComponent,
    ],
    entryComponents: [ TransactionDetailsDialogComponent ],
    providers: [
    ]
})
export class TransactionsModule { }
