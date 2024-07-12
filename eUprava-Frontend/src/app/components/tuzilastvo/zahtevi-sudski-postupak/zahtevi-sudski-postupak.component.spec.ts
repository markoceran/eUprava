import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ZahteviSudskiPostupakComponent } from './zahtevi-sudski-postupak.component';

describe('ZahteviSudskiPostupakComponent', () => {
  let component: ZahteviSudskiPostupakComponent;
  let fixture: ComponentFixture<ZahteviSudskiPostupakComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ZahteviSudskiPostupakComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ZahteviSudskiPostupakComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
