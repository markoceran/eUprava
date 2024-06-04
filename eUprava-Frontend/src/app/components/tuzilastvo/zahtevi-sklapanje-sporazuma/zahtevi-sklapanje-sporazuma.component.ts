import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ZahtevZaSklapanjeSporazuma } from 'src/app/models/zahtevZaSklapanjeSporazuma';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';

@Component({
  selector: 'app-zahtevi-sklapanje-sporazuma',
  templateUrl: './zahtevi-sklapanje-sporazuma.component.html',
  styleUrls: ['./zahtevi-sklapanje-sporazuma.component.css']
})
export class ZahteviSklapanjeSporazumaComponent implements OnInit {

  constructor(private router: Router,private tuzilastvoService:TuzilastvoService) { }

  zahteviZaSklapanjeSporazuma: ZahtevZaSklapanjeSporazuma[] = [];

  ngOnInit(): void {
    this.getZahteviZaSklapanjeSporazuma();
  }

  getZahteviZaSklapanjeSporazuma(): void {
    this.tuzilastvoService.getZahteviZaSklapanjeSporazuma().subscribe(
      (data) => {
        this.zahteviZaSklapanjeSporazuma = data;
      },
      (error) => {
        console.error(error);
      }
    );
  }

}
