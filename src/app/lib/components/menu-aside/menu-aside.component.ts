import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';
import { MarketService, StoreService } from '../../services';
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
  ) {

    this.storeService.ready.subscribe((ready) => {
      this.storeService.getAllAccounts().then((accounts) => {
        this.accounts = accounts;
      })
      this.storeService.getSelectedAccount().then((account) => {
        this.selectedAccount = account;
      })
    });

  }

  public ngOnInit() {
    // TODO
  }

  openSendDialog(): void {

    // get suggested fees
    this.marketService.getSuggestedFees().subscribe((fees) => {

      // open dialog
      const dialogRef = this.dialog.open(SendBurstDialogComponent, {
        width: '600px',
        data: { account: this.selectedAccount, fees: fees }
      });

      dialogRef.afterClosed().subscribe(result => {
        console.log('The dialog was closed', result);
      });
    });

  }

}
