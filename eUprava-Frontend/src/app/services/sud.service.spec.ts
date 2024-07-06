import { TestBed } from '@angular/core/testing';

import { SudService } from './sud.service';

describe('SudService', () => {
  let service: SudService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SudService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
