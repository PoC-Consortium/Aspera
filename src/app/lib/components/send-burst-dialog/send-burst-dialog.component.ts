import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';
import { Account, Transaction } from '../../model';
import { CryptoService, AccountService } from '../../services';


@Component({
  selector: 'app-send-burst-dialog',
  templateUrl: './send-burst-dialog.component.html',
  styleUrls: ['./send-burst-dialog.component.css']
})
export class SendBurstDialogComponent implements OnInit {

  constructor(
    public dialogRef: MatDialogRef<SendBurstDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data,
    public cryptoService: CryptoService,
    public accountService: AccountService) { 
    }

  closeDialog(): void {
    this.dialogRef.close();
  }

  ngOnInit() {
    
  }

  sendBurst(transactionRequest) {
    const { transaction, pin } = transactionRequest;
    let transactionToSend: Transaction = { 
      senderPublicKey: this.data.account.keys.publicKey,
      ...transaction 
    };
    return this.accountService.doTransaction(transactionToSend, this.data.account.keys.signPrivateKey, pin).then((transaction: Transaction) => {
      console.log(transaction);
      this.closeDialog();
    });
  }

}
