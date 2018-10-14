import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';
import { Account } from '../../model';

@Component({
  selector: 'app-send-burst-dialog',
  templateUrl: './send-burst-dialog.component.html',
  styleUrls: ['./send-burst-dialog.component.css']
})
export class SendBurstDialogComponent implements OnInit {

  constructor(
    public dialogRef: MatDialogRef<SendBurstDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Account) { }

  onNoClick(): void {
    this.dialogRef.close();
  }

  ngOnInit() {
  }

}
