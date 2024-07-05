import { Component, OnInit } from '@angular/core';
import { GranicnaPolicijaService } from '../../../../services/granicna-policija.service';
import { Prelaz } from '../../../../models/prelaz'

@Component({
  selector: 'app-prelazi',
  templateUrl: './prelazi.component.html',
  styleUrls: ['./prelazi.component.css']
})
export class PrelaziComponent implements OnInit {
  prelazi: Prelaz[] = [];

  constructor(private granicnaPolicijaService: GranicnaPolicijaService) { }

  ngOnInit(): void {
    this.getPrelazi();
  }

  getPrelazi(): void {
    this.granicnaPolicijaService.getPrelazi().subscribe(prelazi => {
      this.prelazi = prelazi;
    });
  }
}
