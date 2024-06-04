import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KrivicnePrijaveComponent } from './krivicne-prijave.component';

describe('KrivicnePrijaveComponent', () => {
  let component: KrivicnePrijaveComponent;
  let fixture: ComponentFixture<KrivicnePrijaveComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KrivicnePrijaveComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KrivicnePrijaveComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
