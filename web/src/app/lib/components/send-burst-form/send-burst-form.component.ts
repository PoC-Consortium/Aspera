import { Component, OnInit, ViewChild, Input, Output, EventEmitter } from '@angular/core';
import { Transaction, Attachment, SuggestedFees, EncryptedMessage, Message } from '../../model';
import { NgForm } from '@angular/forms';
import { BurstUtil } from '../../util/burst';
import { CryptoService } from '../../services';

@Component({
  selector: 'app-send-burst-form',
  templateUrl: './send-burst-form.component.html',
  styleUrls: ['./send-burst-form.component.css']
})
export class SendBurstFormComponent implements OnInit {
  @ViewChild('sendBurstForm') public sendBurstForm: NgForm;
  @ViewChild('feeNQT') public feeNQT: string;
  @ViewChild('recipientAddress') public recipientAddress: string;
  @ViewChild('amountNQT') public amountNQT: string;
  @ViewChild('message') public message: string;
  @ViewChild('fullHash') public fullHash: string;
  @ViewChild('encrypt') public encrypt: string;
  @ViewChild('pin') public pin: string;
  @ViewChild('deadline') public deadline: string;

  @Input('fees') public fees: SuggestedFees;
  @Input('balance') public balance: number;
  @Input('close') public close: Function;

  @Output() submit = new EventEmitter<any>();
  advanced: boolean = false;
  showMessage: boolean = false;
  burstAddressPattern = BurstUtil.burstAddressPattern;

  constructor() {
  }

  ngOnInit() {
  }

  getTotal() {
    return parseFloat(this.amountNQT) + parseFloat(this.feeNQT) || 0;
  }

  setFee(feeNQT: string) {
    this.feeNQT = this.convertFeeToBurst(feeNQT).toString();
  }

  convertFeeToBurst(feeNQT: string) {
    return BurstUtil.convertStringToNumber(feeNQT);
  }

  onSubmit(event) {
    this.submit.emit({
      transaction: {
        recipientAddress: `BURST-${this.recipientAddress}`,
        amountNQT: parseFloat(this.amountNQT),
        feeNQT: this.feeNQT,
        attachment: this.getMessage(),
        deadline: parseFloat(this.deadline),
        fullHash: this.fullHash,
        type: 1
      },
      pin: this.pin
    });
    event.stopImmediatePropagation();
  }

  getMessage() {
    if (this.message) {
      if (this.encrypt) {
        return {
          encryptedMessage: this.message
        }
      } else {
        return {
          message: this.message,
          type: 'message',
          messageIsText: true
        }
      }
    }
    return null;
  }
}
