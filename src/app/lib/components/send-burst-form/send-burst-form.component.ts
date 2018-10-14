import { Component, OnInit, ViewChild } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';
import { Transaction, Attachment } from '../../model';
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
  @ViewChild('deadline') public deadline: string = "24";

  submitted$: Subject<Transaction>;
  advanced: boolean = false;
  showMessage: boolean = false;

  constructor() { 
    this.submitted$ = new Subject<Transaction>();
  }

  ngOnInit() {
  }

  getSuggestedStandardFee() {
    return "0.04441000"
  }

  getSuggestedCheapFee() {
    return "0.007441000"
  }

  getSuggestedPriorityFee() {
    return "0.1470000"
  }

  getTotal() {
    return parseFloat(this.amountNQT) + parseFloat(this.feeNQT) || 0;
  }

  setFee(feeNQT: string) {
    this.feeNQT = feeNQT;
  }

  onSubmit(values) {
    console.log(values);
    console.log(this.sendBurstForm);
    this.submitted$.next({
      recipientAddress: this.recipientAddress,
      amountNQT: parseFloat(this.amountNQT),
      feeNQT: parseFloat(this.feeNQT),
      attachment: this.getMessage(),
      deadline: parseFloat(this.deadline),
      fullHash: this.fullHash,
      type: 1
    });
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
