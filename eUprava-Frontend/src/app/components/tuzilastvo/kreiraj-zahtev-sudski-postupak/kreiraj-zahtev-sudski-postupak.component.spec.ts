import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KreirajZahtevSudskiPostupakComponent } from './kreiraj-zahtev-sudski-postupak.component';

describe('KreirajZahtevSudskiPostupakComponent', () => {
  let component: KreirajZahtevSudskiPostupakComponent;
  let fixture: ComponentFixture<KreirajZahtevSudskiPostupakComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KreirajZahtevSudskiPostupakComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KreirajZahtevSudskiPostupakComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
