import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { KrivicnaPrijava } from 'src/app/models/krivicnaPrijava';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';
import { KreirajZahtevSklapanjeSporazumaDialogComponent } from '../kreiraj-zahtev-sklapanje-sporazuma-dialog/kreiraj-zahtev-sklapanje-sporazuma-dialog.component';
import { KreirajZahtevSudskiPostupakComponent } from '../kreiraj-zahtev-sudski-postupak/kreiraj-zahtev-sudski-postupak.component';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-krivicne-prijave',
  templateUrl: './krivicne-prijave.component.html',
  styleUrls: ['./krivicne-prijave.component.css']
})
export class KrivicnePrijaveComponent implements OnInit {

  constructor(private authService: AuthService,private tuzilastvoService:TuzilastvoService,public dialog: MatDialog) { }

  krivicnePrijave: KrivicnaPrijava[] = [];
  rolaLogovanogKorisnika: string | null = ""

  ngOnInit(): void {
    this.getKrivicnePrijave();
    this.rolaLogovanogKorisnika = this.authService.extractUserType();
  }

  getKrivicnePrijave(): void {
    this.tuzilastvoService.getKrivicnePrijave().subscribe(
      (data) => {
        this.krivicnePrijave = data;
      },
      (error) => {
        console.error(error);
      }
    );
  }

  openKreirajZahtevZaSklapanjeSporazumaDialog(krivicnaPrijavaId:any): void {
    const dialogRef = this.dialog.open(KreirajZahtevSklapanjeSporazumaDialogComponent, {
      width: '400px',
      data: { krivicnaPrijavaId: krivicnaPrijavaId }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        console.log('The dialog was closed with result:', result);
        // Optionally handle the result here
      }
    });
  }

  openKreirajZahtevZaSudskiPostupakDialog(krivicnaPrijavaId:any): void {
    const dialogRef = this.dialog.open(KreirajZahtevSudskiPostupakComponent, {
      width: '400px',
      data: { krivicnaPrijavaId: krivicnaPrijavaId }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        console.log('The dialog was closed with result:', result);
        // Optionally handle the result here
      }
    });
  }
}
