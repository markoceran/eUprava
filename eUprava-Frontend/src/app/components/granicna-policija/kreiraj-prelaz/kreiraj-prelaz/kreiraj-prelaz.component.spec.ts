import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KreirajPrelazComponent } from './kreiraj-prelaz.component';

describe('KreirajPrelazComponent', () => {
  let component: KreirajPrelazComponent;
  let fixture: ComponentFixture<KreirajPrelazComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KreirajPrelazComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KreirajPrelazComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
