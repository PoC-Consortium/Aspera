import { Component, OnInit, ViewChild, Input, Output, EventEmitter } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';
import { Transaction, Attachment, SuggestedFees } from '../../model';
import { NgForm } from '@angular/forms';
import { BurstUtil } from '../../util/burst';
import { AccountService } from '../../services';

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
      : calculateMultiOutTotal + parseFloat(this.feeNQT) || 0;
  }

  setFee(feeNQT: string) {
    this.feeNQT = this.convertFeeToBurst(feeNQT).toString();
  }

  convertFeeToBurst(feeNQT: string) {
    return BurstUtil.convertStringToNumber(feeNQT);
  }

  getMultiOutString() {
    if (this.sameAmount) {
      return this.recipients.map((recipient) => 
        `${BurstUtil.decode(`BURST-${recipient.address}`)}`)
        .join(';');
    } else {
      return this.recipients.map((recipient) => 
        `${BurstUtil.decode(`BURST-${recipient.address}`)}:${BurstUtil.convertNumberToString(parseFloat(recipient.amountNQT))}`)
        .join(';');
    }
  }
  
  onSubmit(event) {
    const multiOutString = this.getMultiOutString();

    this.submit.emit({
      transaction: {
        recipients: multiOutString,
        feeNQT: BurstUtil.convertNumberToString(parseFloat(this.feeNQT)),
        deadline: "1440",
        amountNQT: BurstUtil.convertNumberToString(this.getTotal())
      },
      pin: this.pin,
      sameAmount: this.sameAmount
    });
    event.stopImmediatePropagation();
  }
}

