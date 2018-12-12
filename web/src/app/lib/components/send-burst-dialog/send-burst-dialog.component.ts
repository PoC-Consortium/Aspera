import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';
import { Account, Transaction, EncryptedMessage } from '../../model';
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

  async sendBurst(transactionRequest) {
    const { transaction, pin } = transactionRequest;
    let transactionToSend: Transaction = { 
      senderPublicKey: this.data.account.keys.publicKey,
      ...transaction 
    };

    // todo: move to service
    if (transactionToSend.attachment && transactionToSend.attachment.encryptedMessage) {
      const recipientPublicKey = await this.accountService.getAccountPublicKey(transaction.recipientAddress);
      const encryptedMessage = await this.cryptoService.encryptMessage(transactionToSend.attachment.encryptedMessage, 
        this.data.account.keys.agreementPrivateKey, this.accountService.hashPinEncryption(pin), recipientPublicKey);
      transactionToSend.attachment = new EncryptedMessage({
        data: encryptedMessage.m,
        nonce: encryptedMessage.n,
        isText: true
      })
    }

    return this.accountService.doTransaction(transactionToSend, this.data.account.keys.signPrivateKey, pin).then((transaction: Transaction) => {
      console.log(transaction);
      this.closeDialog();
    });
  }


  sendMessage(transactionRequest) {
    const { transaction, pin } = transactionRequest;
    let transactionToSend: Transaction = { 
      senderPublicKey: this.data.account.keys.publicKey,
      ...transaction 
    };

    return this.accountService.sendMessage(transactionToSend, this.data.account.keys.signPrivateKey, pin).then((transaction: Transaction) => {
      console.log(transaction);
      this.closeDialog();
    });
  }

  sendBurstMultiOut(transactionRequest) {
    const { transaction, pin, sameAmount } = transactionRequest;
    let transactionToSend: Transaction = { 
      senderPublicKey: this.data.account.keys.publicKey,
      ...transaction 
    };
    return this.accountService.doMultiOutTransaction(transactionToSend, this.data.account.keys.signPrivateKey, pin, sameAmount).then((transaction: Transaction) => {
      console.log(transaction);
      this.closeDialog();
    });
  }

}
