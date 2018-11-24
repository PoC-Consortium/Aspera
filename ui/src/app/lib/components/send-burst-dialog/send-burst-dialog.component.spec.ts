import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SendBurstDialogComponent } from './send-burst-dialog.component';

describe('SendBurstDialogComponent', () => {
  let component: SendBurstDialogComponent;
  let fixture: ComponentFixture<SendBurstDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SendBurstDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendBurstDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
