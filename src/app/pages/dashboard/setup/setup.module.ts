import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { SharedModule } from '../../../lib/shared.module';

import { SetupRouting } from './setup.routing';
import { SetupComponent } from './setup.component';
import { NodeSetupComponent } from './node/node.component';
import { NodeSetupAddressComponent } from './node/address/address.component';
import { AccountNewComponent } from './account/account.component';
import { AccountCreateComponent } from './account/create/create.component';
import { AccountCreatePinComponent } from './account/create/pin/pin.component';
import { AccountCreateRecordComponent } from './account/create/record/record.component';
import { AccountCreateSeedComponent } from './account/create/seed/seed.component';

import { SetupService } from './setup.service';
import { CreateService } from './account/create/create.service';
import { NodeService } from './node/node.service';

import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material';
import { MatListModule } from '@angular/material/list';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
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
        MatListModule,
        MatProgressBarModule,
        MatProgressSpinnerModule,
        MatSlideToggleModule,
        MatStepperModule,
        ReactiveFormsModule,
        SetupRouting,
        SharedModule,
        FormsModule
    ],
    declarations: [
        SetupComponent,
        AccountCreateComponent,
        AccountCreatePinComponent,
        AccountCreateRecordComponent,
        AccountCreateSeedComponent,
        AccountNewComponent,
        NodeSetupComponent,
        NodeSetupAddressComponent
    ],
    providers: [
        CreateService,
        NodeService,
        SetupService
    ],
    exports: [
        AccountCreateComponent
    ]
})
export class SetupModule { }
