import { Component, OnInit } from '@angular/core';
import { GranicnaPolicijaService } from '../../../../services/granicna-policija.service';
import { SumnjivoLice } from '../../../../models/sumnjivoLice';

@Component({
  selector: 'app-sumnjiva-lica',
  templateUrl: './sumnjiva-lica.component.html',
  styleUrls: ['./sumnjiva-lica.component.css']
})
export class SumnjivaLicaComponent implements OnInit {
  sumnjivaLica: SumnjivoLice[] = [];

  constructor(private granicnaPolicijaService: GranicnaPolicijaService) { }

  ngOnInit(): void {
    this.getSumnjivaLica();
  }

  getSumnjivaLica(): void {
    this.granicnaPolicijaService.getSumnjivaLica().subscribe(sumnjivaLica => {
      this.sumnjivaLica = sumnjivaLica;
    });
  }
}
