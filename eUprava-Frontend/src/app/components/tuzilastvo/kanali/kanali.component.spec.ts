import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KanaliComponent } from './kanali.component';

describe('KanaliComponent', () => {
  let component: KanaliComponent;
  let fixture: ComponentFixture<KanaliComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KanaliComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KanaliComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
