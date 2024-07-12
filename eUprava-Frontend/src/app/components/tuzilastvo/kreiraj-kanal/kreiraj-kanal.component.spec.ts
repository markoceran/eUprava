import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KreirajKanalComponent } from './kreiraj-kanal.component';

describe('KreirajKanalComponent', () => {
  let component: KreirajKanalComponent;
  let fixture: ComponentFixture<KreirajKanalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KreirajKanalComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KreirajKanalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
