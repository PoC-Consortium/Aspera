import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import {} from 'jasmine';
import { SendBurstDialogComponent } from './send-burst-dialog.component';
import { ElementRef } from '@angular/core';

export class MockElementRef extends ElementRef {nativeElement = {};}

xdescribe('SendBurstDialogComponent', () => {
  let component: SendBurstDialogComponent;
  let fixture: ComponentFixture<SendBurstDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SendBurstDialogComponent ],
      providers: [
        { provide: ElementRef, useClass: MockElementRef },
      ]
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
