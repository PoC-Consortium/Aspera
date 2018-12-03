import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { CreatePassiveAccountComponent } from './create-passive.component';
import { ElementRef } from '@angular/core';
import { MockElementRef } from 'src/app/lib/components/send-burst-form/burst-input-validator.directive.spec';


describe('CreatePassiveAccountComponent', () => {
  let component: CreatePassiveAccountComponent;
  let fixture: ComponentFixture<CreatePassiveAccountComponent>;

  beforeEach(async(() => { 
    TestBed.configureTestingModule({
      declarations: [ CreatePassiveAccountComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreatePassiveAccountComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy(); 
  });
});
