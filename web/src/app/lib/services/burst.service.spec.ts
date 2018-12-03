import { TestBed, inject } from '@angular/core/testing';

import { BurstService } from './burst.service';

xdescribe('BurstService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [BurstService]
    });
  });

  it('should be created', inject([BurstService], (service: BurstService) => {
    expect(service).toBeTruthy();
  }));
});
