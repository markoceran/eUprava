import { Component, OnDestroy, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { ZahtevZaSklapanjeSporazuma } from 'src/app/models/zahtevZaSklapanjeSporazuma';
import { AuthService } from 'src/app/services/auth.service';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';

@Component({
  selector: 'app-zahtevi-sklapanje-sporazuma',
  templateUrl: './zahtevi-sklapanje-sporazuma.component.html',
  styleUrls: ['./zahtevi-sklapanje-sporazuma.component.css']
})
export class ZahteviSklapanjeSporazumaComponent implements OnInit {

  zahteviZaSklapanjeSporazuma: ZahtevZaSklapanjeSporazuma[] = [];
  role!: string | null;

  constructor(
    private router: Router,
    private tuzilastvoService: TuzilastvoService,
    private authService: AuthService
  ) { }

  ngOnInit(): void {
      if (this.authService.isLoggedIn()) {
        this.role = this.authService.extractUserType();
        console.log(this.role)
        if (this.role != null && this.role === 'Gradjanin') {
          const id = this.authService.getUserIdFromToken();
          if (id) {
            this.getZahteviZaSklapanjeSporazumaPoGradjaninu(id);
          } else {
            console.error('Id korisnika nije pronadjen.');
          }
        } else if (this.role != null && this.role === 'Tuzioc') {
          this.getZahteviZaSklapanjeSporazuma();
        }
      } else {
        this.zahteviZaSklapanjeSporazuma = [];
      }
  }  
  

  getZahteviZaSklapanjeSporazuma(): void {
    this.tuzilastvoService.getZahteviZaSklapanjeSporazuma().subscribe(
      (data) => {
        this.zahteviZaSklapanjeSporazuma = data;
      },
      (error) => {
        console.error('Greska prilikom dobavljanja:', error);
      }
    );
  }

  getZahteviZaSklapanjeSporazumaPoGradjaninu(id: any): void {
    this.tuzilastvoService.getZahteviZaSklapanjeSporazumaPoGradjaninu(id).subscribe(
      (data) => {
        this.zahteviZaSklapanjeSporazuma = data;
      },
      (error) => {
        console.error(`Greska prilikom dobavljanja:`, error);
      }
    );
  }

}
