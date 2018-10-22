import { Component, OnInit } from '@angular/core';
import { CreateService } from '../create.service';
import { BurstUtil } from '../../../../../lib/util';
import { NotifierService } from 'angular-notifier';
import { Router } from '@angular/router';

@Component({
  selector: 'app-create-passive',
  templateUrl: './create-passive.component.html',
  styleUrls: ['./create-passive.component.css']
})
export class CreatePassiveAccountComponent implements OnInit {

  burstAddressPattern = BurstUtil.burstAddressPattern;

  constructor(private createService: CreateService,
              private notificationService: NotifierService,
              private router: Router) { }

  ngOnInit() {
  }

  public submit(address: string) {
      this.createService.setAddress(address);
      this.createService.createPassiveAccount().then((success) => {
        this.notificationService.notify('success', `Account added: ${address}`);
        this.router.navigate(['/']);
      },
      (error) => {
        this.notificationService.notify('error', error.toString());
      });
  }

}
