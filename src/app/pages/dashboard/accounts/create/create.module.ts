import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { SharedModule } from '../../../../lib/shared.module';

import { CreateComponent } from './create.component';
import { CreateRouting } from './create.routing';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material';
import { MatStepperModule } from '@angular/material/stepper';

@NgModule({
    imports: [
        CreateRouting,
        CommonModule,
        ReactiveFormsModule,
        SharedModule,
        MatCardModule,
        MatFormFieldModule,
        MatGridListModule,
        MatIconModule,
        MatInputModule,
        MatStepperModule
    ],
    declarations: [
        CreateComponent
    ],
    providers: [
        //CreateService
    ]
})
export class AccountsCreateModule { }
