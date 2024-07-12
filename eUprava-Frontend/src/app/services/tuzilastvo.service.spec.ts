import { TestBed } from '@angular/core/testing';

import { TuzilastvoService } from './tuzilastvo.service';

describe('TuzilastvoService', () => {
  let service: TuzilastvoService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TuzilastvoService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
