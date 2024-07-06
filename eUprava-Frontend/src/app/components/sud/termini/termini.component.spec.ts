import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TerminiComponent } from './termini.component';

describe('TerminiComponent', () => {
  let component: TerminiComponent;
  let fixture: ComponentFixture<TerminiComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TerminiComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TerminiComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
