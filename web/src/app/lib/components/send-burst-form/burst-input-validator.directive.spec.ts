import { BurstInputValidatorDirective } from './burst-input-validator.directive';
import { TestBed, async } from '@angular/core/testing';
import { ElementRef } from '@angular/core';
import { NgControl, FormControl } from '@angular/forms';

export class MockElementRef extends ElementRef {nativeElement = {};}

import { FormsModule } from '@angular/forms';

beforeEach(async(() => {

  TestBed.configureTestingModule({
    imports: [ FormsModule ],
    providers: [
      { provide: ElementRef, useClass: MockElementRef }
    ]
  }).compileComponents();
}));

xdescribe('BurstInputValidatorDirective', () => {
  let mockEl = new MockElementRef({});

  it('should create an instance', () => {
    // const directive = new BurstInputValidatorDirective(mockEl, new FormControl());
    // expect(directive).toBeTruthy();
  });
});
