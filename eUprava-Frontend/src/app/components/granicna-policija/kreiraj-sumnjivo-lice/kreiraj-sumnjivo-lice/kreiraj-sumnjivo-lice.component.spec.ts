import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KreirajSumnjivoLiceComponent } from './kreiraj-sumnjivo-lice.component';

describe('KreirajSumnjivoLiceComponent', () => {
  let component: KreirajSumnjivoLiceComponent;
  let fixture: ComponentFixture<KreirajSumnjivoLiceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KreirajSumnjivoLiceComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KreirajSumnjivoLiceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
