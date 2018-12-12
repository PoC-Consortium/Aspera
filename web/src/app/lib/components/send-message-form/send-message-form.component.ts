import { Component, OnInit, ViewChild, Input, Output, EventEmitter } from '@angular/core';
import { Transaction, Attachment, SuggestedFees, EncryptedMessage, Message } from '../../model';
import { NgForm } from '@angular/forms';
import { BurstUtil } from '../../util/burst';
import { CryptoService } from '../../services';

@Component({
  selector: 'app-send-message-form',
  templateUrl: './send-message-form.component.html',
  styleUrls: ['./send-message-form.component.css']
})
export class SendMessageFormComponent implements OnInit {
  @ViewChild('sendMessageForm') public sendMessageForm: NgForm;
  @ViewChild('feeNQT') public feeNQT: string;
  @ViewChild('recipientAddress') public recipientAddress: string;
  @ViewChild('message') public message: string;
  @ViewChild('fullHash') public fullHash: string;
  @ViewChild('encrypt') public encrypt: string;
  @ViewChild('pin') public pin: string;
  @ViewChild('deadline') public deadline: string = "24";

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
        feeNQT: this.feeNQT,
        attachment: this.getMessage(),
        deadline: parseFloat(this.deadline) * 60,
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