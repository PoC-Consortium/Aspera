import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { EncryptedMessage, Message, Account } from 'src/app/lib/model';
import { CryptoService, BurstService, StoreService } from 'src/app/lib/services';

type TransactionDetailsCellValue = string | Message | EncryptedMessage | number;
type TransactionDetailsCellValueMap = [string, TransactionDetailsCellValue];

@Component({
  selector: 'app-transaction-details-dialog',
  templateUrl: './transaction-details-dialog.component.html',
  styleUrls: ['./transaction-details-dialog.component.css']
})
export class TransactionDetailsDialogComponent implements OnInit {

  infoColumns = ['key', 'value'];
  infoRows = ['transactionType', 'attachment', 'amountNQT', 'feeNQT', 'senderAddress', 'recipientAddress'];
  infoData: Map<string, TransactionDetailsCellValue>;
  detailsData: Map<string, TransactionDetailsCellValue>;
  account: Account;

  constructor(
    public dialogRef: MatDialogRef<TransactionDetailsDialogComponent>,
    private burstService: BurstService,
    private storeService: StoreService,
    @Inject(MAT_DIALOG_DATA) public data) { 
      
      this.storeService.getSelectedAccount()
        .then((account) => {
            this.account = account;
            data.transaction.transactionType = this.burstService.getTransactionNameFromType(this.data.transaction, this.account);
            const transactionDetails = Object.keys(data.transaction).map((key:string): TransactionDetailsCellValueMap => [ key, data.transaction[key]]);
            this.detailsData = new Map(transactionDetails);
            this.infoData = new Map(transactionDetails.filter((row) => this.infoRows.indexOf(row[0]) > -1));
        });
  }

  closeDialog(): void {
    this.dialogRef.close();
  }

  public getInfoData(): TransactionDetailsCellValueMap[] {
    return Array.from(this.infoData.entries());
  } 

  public getDetailsData(): TransactionDetailsCellValueMap[] {
    return Array.from(this.detailsData.entries());
  } 

  ngOnInit() {
     
  }

}
