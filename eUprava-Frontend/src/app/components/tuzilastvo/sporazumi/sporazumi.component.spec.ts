import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SporazumiComponent } from './sporazumi.component';

describe('SporazumiComponent', () => {
  let component: SporazumiComponent;
  let fixture: ComponentFixture<SporazumiComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SporazumiComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SporazumiComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
