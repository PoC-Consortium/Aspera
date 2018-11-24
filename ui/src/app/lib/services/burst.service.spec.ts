import { TestBed, inject } from '@angular/core/testing';

import { BurstService } from './burst.service';

describe('BurstService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [BurstService]
    });
  });

  it('should be created', inject([BurstService], (service: BurstService) => {
    expect(service).toBeTruthy();
  }));
});
