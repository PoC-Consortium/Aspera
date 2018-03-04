import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { SharedModule } from '../../../lib/shared.module';

import { SetupRouting } from './setup.routing';
import { SetupComponent } from './setup.component';
import { AccountNewComponent } from './account/account.component';
import { AccountCreateComponent } from './account/create/create.component';
import { AccountCreateRecordComponent } from './account/create/record/record.component';
import { AccountCreateSeedComponent } from './account/create/seed/seed.component';

import { CreateService } from './account/create/create.service';

import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatStepperModule } from '@angular/material/stepper';

@NgModule({
    imports: [
        CommonModule,
        MatButtonModule,
        MatCardModule,
        MatChipsModule,
        MatDividerModule,
        MatExpansionModule,
        MatFormFieldModule,
        MatGridListModule,
        MatIconModule,
        MatInputModule,
        MatProgressBarModule,
        MatProgressSpinnerModule,
        MatStepperModule,
        ReactiveFormsModule,
        SetupRouting,
        SharedModule
    ],
    declarations: [
        SetupComponent,
        AccountCreateComponent,
        AccountCreateRecordComponent,
        AccountCreateSeedComponent,
        AccountNewComponent
    ],
    providers: [
        CreateService
    ]
})
export class SetupModule { }
