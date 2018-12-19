import { Component, OnInit, Input } from '@angular/core';
import { Attachment, EncryptedMessage, Message } from 'src/app/lib/model';
import { BurstService } from 'src/app/lib/services';
import { BurstUtil } from 'src/app/lib/util';

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
    switch(this.value && this.value.constructor) {
      case Message: {
        this.valueType = 'Message';
        break;
      }
      case EncryptedMessage: {
        this.valueType = 'EncryptedMessage';
        break;
      }
      case String:
        if (BurstUtil.isBurstcoinAddress(this.value as string)) {
          this.valueType = 'BurstAddress';
        }
    }
  }


  public showPinPrompt() {
    console.log('show pin prompt');
    // return this.cryptoService.decryptMessage(this.value.data, this.value.nonce, encryptedPrivateKey, pinHash, )
  }

  public showAccountDialog(address: string) {
    console.log('show account dialog', address);
  }
}
 