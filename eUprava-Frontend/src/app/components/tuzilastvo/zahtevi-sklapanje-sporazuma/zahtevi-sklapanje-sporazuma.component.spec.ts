import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ZahteviSklapanjeSporazumaComponent } from './zahtevi-sklapanje-sporazuma.component';

describe('ZahteviSklapanjeSporazumaComponent', () => {
  let component: ZahteviSklapanjeSporazumaComponent;
  let fixture: ComponentFixture<ZahteviSklapanjeSporazumaComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ZahteviSklapanjeSporazumaComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ZahteviSklapanjeSporazumaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
