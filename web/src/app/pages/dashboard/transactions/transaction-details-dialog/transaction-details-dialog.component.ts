import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { EncryptedMessage, Message } from 'src/app/lib/model';
import { CryptoService, BurstService } from 'src/app/lib/services';


@Component({
  selector: 'app-transaction-details-dialog',
  templateUrl: './transaction-details-dialog.component.html',
  styleUrls: ['./transaction-details-dialog.component.css']
})
export class TransactionDetailsDialogComponent implements OnInit {

  infoColumns = ['key', 'value'];
  infoRows = ['transactionType', 'attachment', 'amountNQT', 'feeNQT', 'senderAddress', 'recipientAddress'];
  infoData: Map<string, string | Message | EncryptedMessage | number>;
  detailsData: Map<string, string | Message | EncryptedMessage | number>;
  encryptedMessage = false;
  unencryptedMessage = false;

  constructor(
    public dialogRef: MatDialogRef<TransactionDetailsDialogComponent>,
    private burstService: BurstService,
    @Inject(MAT_DIALOG_DATA) public data) { 
      data.transaction.transactionType = this.burstService.getTransactionNameFromType(this.data.transaction);
      const transactionDetails = Object.keys(data.transaction).map((key:string): [string, string | Message | EncryptedMessage | number] => [ key, data.transaction[key]]);
      this.detailsData = new Map(transactionDetails);
      this.infoData = new Map(transactionDetails.filter((row) => this.infoRows.indexOf(row[0]) > -1));
  }

  closeDialog(): void {
    this.dialogRef.close();
  }

  public getInfoData(): [string, string | Message | EncryptedMessage | number][] {
    return Array.from(this.infoData.entries());
  } 

  public getDetailsData(): [string, string | Message | EncryptedMessage | number][] {
    return Array.from(this.detailsData.entries());
  } 

  ngOnInit() {
     
  }

}
