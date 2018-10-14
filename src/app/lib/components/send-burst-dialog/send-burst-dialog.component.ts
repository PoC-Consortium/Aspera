import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';
import { Account } from '../../model';
import { CryptoService } from '../../services';

@Component({
  selector: 'app-send-burst-dialog',
  templateUrl: './send-burst-dialog.component.html',
  styleUrls: ['./send-burst-dialog.component.css']
})
export class SendBurstDialogComponent implements OnInit {

  constructor(
    public dialogRef: MatDialogRef<SendBurstDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data,
    public cryptoService: CryptoService) { 
    }

  onNoClick(): void {
    this.dialogRef.close();
  }

  ngOnInit() {
    
  }

}
