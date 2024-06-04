import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KreirajZahtevSklapanjeSporazumaDialogComponent } from './kreiraj-zahtev-sklapanje-sporazuma-dialog.component';

describe('KreirajZahtevSklapanjeSporazumaDialogComponent', () => {
  let component: KreirajZahtevSklapanjeSporazumaDialogComponent;
  let fixture: ComponentFixture<KreirajZahtevSklapanjeSporazumaDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KreirajZahtevSklapanjeSporazumaDialogComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KreirajZahtevSklapanjeSporazumaDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
