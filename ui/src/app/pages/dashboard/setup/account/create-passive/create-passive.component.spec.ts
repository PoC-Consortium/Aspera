import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CreatePassiveAccountComponent } from './create-passive.component';

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
