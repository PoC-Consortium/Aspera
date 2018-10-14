import { Component, OnInit, ViewChild, Input, Output, EventEmitter } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';
import { Transaction, Attachment, SuggestedFees } from '../../model';
import { NgForm } from '@angular/forms';

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

  @Output() submit = new EventEmitter<any>();
  advanced: boolean = false;
  showMessage: boolean = false;

  constructor() {
  }

  ngOnInit() {
  }

  getTotal() {
    return parseFloat(this.amountNQT) + parseFloat(this.feeNQT) || 0;
  }

  setFee(feeNQT: string) {
    this.feeNQT = feeNQT;
  }

  onSubmit(event) {
    this.submit.emit({
      transaction: {
        recipientAddress: this.recipientAddress,
        amountNQT: parseFloat(this.amountNQT),
        feeNQT: parseFloat(this.feeNQT),
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
          message: this.message
        }
      }
    }
    return null;
  }

}
