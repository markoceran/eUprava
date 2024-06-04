import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ZahtevZaSudskiPostupak } from 'src/app/models/zahtevZaSudskiPostupak';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';

@Component({
  selector: 'app-zahtevi-sudski-postupak',
  templateUrl: './zahtevi-sudski-postupak.component.html',
  styleUrls: ['./zahtevi-sudski-postupak.component.css']
})
export class ZahteviSudskiPostupakComponent implements OnInit {

  constructor(private router: Router,private tuzilastvoService:TuzilastvoService) { }

  zahteviZaSudskiPostupak: ZahtevZaSudskiPostupak[] = [];

  ngOnInit(): void {
    this.getZahteviZaSudskiPostupak();
  }

  getZahteviZaSudskiPostupak(): void {
    this.tuzilastvoService.getZahteviZaSudskiPostupak().subscribe(
      (data) => {
        this.zahteviZaSudskiPostupak = data;
      },
      (error) => {
        console.error(error);
      }
    );
  }

}
