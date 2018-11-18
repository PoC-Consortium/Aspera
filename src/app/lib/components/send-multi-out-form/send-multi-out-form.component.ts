import { Component, OnInit, ViewChild, Input, Output, EventEmitter } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';
import { Transaction, Attachment, SuggestedFees } from '../../model';
import { NgForm } from '@angular/forms';
import { BurstUtil } from '../../util/burst';

interface Recipient {
  address: string,
  amountNQT: string
}

@Component({
  selector: 'app-send-multi-out-form',
  templateUrl: './send-multi-out-form.component.html',
  styleUrls: ['./send-multi-out-form.component.css']
})
export class SendMultiOutFormComponent implements OnInit {
  @ViewChild('sendBurstForm') public sendBurstForm: NgForm;
  @ViewChild('feeNQT') public feeNQT: string;
  @ViewChild('amountNQT') public amountNQT: string;
  @ViewChild('recipients') public recipients: Recipient[];
  @ViewChild('message') public message: string;
  @ViewChild('fullHash') public fullHash: string;
  @ViewChild('pin') public pin: string;

  @Input('fees') public fees: SuggestedFees;
  @Input('balance') public balance: number;
  @Input('close') public close: Function;

  @Output() submit = new EventEmitter<any>();
  sameAmount: boolean = false;
  burstAddressPattern = BurstUtil.burstAddressPattern;

  constructor() {
  }

  ngOnInit() {
    this.recipients = [this.createRecipient(), this.createRecipient()];
  }

  createRecipient() {
    return { 
      amountNQT:'', 
      address:''
    };
  }

  toggleSameAmount() {
    this.sameAmount = !this.sameAmount;
    console.log(this.sameAmount);
  }

  addRecipient() {
    this.recipients.push(this.createRecipient());
  }

  getTotal() {
    const calculateMultiOutTotal = this.recipients.map((recipient) => {
      return parseFloat(recipient.amountNQT) || 0;
    }).reduce((acc, curr) => acc + curr, 0);

    return this.sameAmount ? parseFloat(this.amountNQT) + parseFloat(this.feeNQT) || 0
      : calculateMultiOutTotal || 0;
  }

  setFee(feeNQT: string) {
    this.feeNQT = this.convertFeeToBurst(feeNQT).toString();
  }

  convertFeeToBurst(feeNQT: string) {
    return parseFloat(feeNQT)/100000000;
  }

  onSubmit(event) {
    // this.submit.emit({
    //   transaction: {
    //     // recipientAddress: `BURST-${this.recipientAddress}`,
    //     amountNQT: parseFloat(this.amountNQT),
    //     feeNQT: this.feeNQT,
    //     attachment: this.getMessage(),
    //     deadline: parseFloat(this.deadline),
    //     fullHash: this.fullHash,
    //     type: 1
    //   },
    //   pin: this.pin
    // });
    event.stopImmediatePropagation();
  }
}

