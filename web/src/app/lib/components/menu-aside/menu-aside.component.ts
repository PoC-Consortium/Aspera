import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';
import { MarketService, StoreService, AccountService } from '../../services';
import { MatDialog } from '@angular/material';
import { Account } from '../../model';
import { SendBurstDialogComponent } from '../send-burst-dialog';

@Component({
  selector: 'app-menu-aside',
  styleUrls: ['./menu-aside.component.scss'],
  templateUrl: './menu-aside.component.html'
})
export class MenuAsideComponent implements OnInit {
  private currentUrl: string;

  @Input() private links: Array<any> = [];

  public selectedAccount: Account;
  public accounts: Account[];

  constructor(
    private marketService: MarketService,
    private storeService: StoreService,
    public dialog: MatDialog,
    private accountService: AccountService
  ) {
  }

  ngOnInit() {
    this.accountService.currentAccount.subscribe((account) => {
        this.getSelectedAccounts();
        this.selectedAccount = account;
    })
  }

  private getSelectedAccounts() {
      this.storeService.getAllAccounts().then((accounts) => {
          this.accounts = accounts;
      });
  }

  public selectAccount(account: Account) {
      this.accountService.selectAccount(account);
  }

}
