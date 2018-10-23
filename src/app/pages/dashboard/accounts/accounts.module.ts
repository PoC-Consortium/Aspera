import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SharedModule } from '../../../lib/shared.module';

import { AccountsComponent } from './accounts.component';
import { AccountsRouting } from './accounts.routing';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule, MatSortModule, MatCheckboxModule, MatButtonModule, MatDialogModule } from '@angular/material';
import { MatTableModule } from '@angular/material/table';
import { DeleteAccountDialogComponent } from './delete-account-dialog/delete-account-dialog.component';

@NgModule({
    imports: [
        AccountsRouting,
        CommonModule,
        FormsModule,
        SharedModule,
        MatCardModule,
        MatFormFieldModule,
        MatGridListModule,
        MatIconModule,
        MatInputModule,
        MatSortModule,
        MatTableModule,
        MatCheckboxModule,
        MatButtonModule,
        MatDialogModule
    ],
    declarations: [
        AccountsComponent,
        DeleteAccountDialogComponent
    ],
    providers: [
    ],
    entryComponents: [ DeleteAccountDialogComponent ]
})
export class AccountsModule { }
