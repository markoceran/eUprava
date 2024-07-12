import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KreirajKrivicnuPrijavuComponent } from './kreiraj-krivicnu-prijavu.component';

describe('KreirajKrivicnuPrijavuComponent', () => {
  let component: KreirajKrivicnuPrijavuComponent;
  let fixture: ComponentFixture<KreirajKrivicnuPrijavuComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KreirajKrivicnuPrijavuComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KreirajKrivicnuPrijavuComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
