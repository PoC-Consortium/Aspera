import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormControl, NgModel, NgForm } from '@angular/forms';
import { Credentials } from '../models/credentials';

@Component({
  selector: 'bc-login-form',
  templateUrl: './login-form.component.html',
  styleUrls: ['./login-form.component.scss'],
})
export class LoginFormComponent implements OnInit {

  passphrase: string;

  @Input() pending: boolean;

  @Input() errorMessage: string | null;

  @Output() submitted = new EventEmitter<Credentials>();

  constructor() {}

  ngOnInit() {}

  submit() {
    this.submitted.emit({ passphrase: this.passphrase });
    this.passphrase = '';
  }
}
