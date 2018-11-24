import { Component, OnInit } from '@angular/core';
import { Store, select } from '@ngrx/store';
import { AccountService } from '../../lib/services';
import { CreateService } from '../dashboard/setup/account/create.service';
import { NotifierService } from 'angular-notifier';
import { Router } from '@angular/router';

@Component({
  selector: 'bc-login-page',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent implements OnInit {

  constructor(private accountService: AccountService,
    private createService: CreateService,
    private notificationService: NotifierService,
    private router: Router) {}

  method: string;

  ngOnInit() {
    this.method = 'active';
  }
}
