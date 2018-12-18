import { Component, OnInit, Input } from '@angular/core';
import { Attachment, EncryptedMessage, Message } from 'src/app/lib/model';
import { BurstService } from 'src/app/lib/services';

@Component({
  selector: 'app-transaction-row-value-cell',
  templateUrl: './transaction-row-value-cell.component.html',
  styleUrls: ['./transaction-row-value-cell.component.css']
})
export class TransactionRowValueCellComponent implements OnInit {

  @Input('value') value: string | Attachment | number;
  @Input('key') key: string;
  valueType: string = "string";

  constructor() { }

  ngOnInit() {
    console.log(this.value);
    switch(this.value && this.value.constructor) {
      case Message: {
        this.valueType = 'Message';
        break;
      }
      case EncryptedMessage: {
        this.valueType = 'EncryptedMessage';
        break;
      }
    }
  }

  public showPinPrompt() {
    console.log('show pin prompt');
    // return this.cryptoService.decryptMessage(this.value.data, this.value.nonce, encryptedPrivateKey, pinHash, )
  }
}
 