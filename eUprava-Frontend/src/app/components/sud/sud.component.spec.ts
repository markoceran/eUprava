import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SudComponent } from './sud.component';

describe('SudComponent', () => {
  let component: SudComponent;
  let fixture: ComponentFixture<SudComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SudComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SudComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
