import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import { LoginPageComponent } from './containers/login-page.component';
import { LoginFormComponent } from './components/login-form.component';
import { AccountsListComponent } from './components/accounts-list/accounts-list.component';
import { LogoutConfirmationDialogComponent } from './components/logout-confirmation-dialog.component';

import { AuthEffects } from './effects/auth.effects';
import { reducers } from './reducers';
import { AuthRoutingModule } from './auth-routing.module';
import { MatCardModule, MatFormFieldModule, MatInputModule, MatButtonModule, MatDialogModule, MatIconModule } from '@angular/material';
import { AccountsListEffects } from './effects/accounts-list.effects';

export const COMPONENTS = [
  LoginPageComponent,
  LoginFormComponent,
  LogoutConfirmationDialogComponent,
  AccountsListComponent
];

@NgModule({
  imports: [
    CommonModule,
    ReactiveFormsModule,
    AuthRoutingModule,
    FormsModule,
    MatIconModule,
    MatInputModule,
    MatButtonModule,
    MatFormFieldModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatDialogModule,
    StoreModule.forFeature('auth', reducers),
    EffectsModule.forFeature([AuthEffects, AccountsListEffects]),
  ],
  declarations: COMPONENTS,
  entryComponents: [LogoutConfirmationDialogComponent],
})
export class AuthModule {}
