import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PrelaziComponent } from './prelazi.component';

describe('PrelaziComponent', () => {
  let component: PrelaziComponent;
  let fixture: ComponentFixture<PrelaziComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PrelaziComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PrelaziComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
