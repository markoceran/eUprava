import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SumnjivaLicaComponent } from './sumnjiva-lica.component';

describe('SumnjivaLicaComponent', () => {
  let component: SumnjivaLicaComponent;
  let fixture: ComponentFixture<SumnjivaLicaComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SumnjivaLicaComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SumnjivaLicaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
