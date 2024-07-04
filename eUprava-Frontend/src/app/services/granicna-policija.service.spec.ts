import { TestBed } from '@angular/core/testing';

import { GranicnaPolicijaService } from './granicna-policija.service';

describe('GranicnaPolicijaService', () => {
  let service: GranicnaPolicijaService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GranicnaPolicijaService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
