import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { KrivicnaPrijava } from 'src/app/models/krivicnaPrijava';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';
import { KreirajZahtevSklapanjeSporazumaDialogComponent } from '../kreiraj-zahtev-sklapanje-sporazuma-dialog/kreiraj-zahtev-sklapanje-sporazuma-dialog.component';

@Component({
  selector: 'app-krivicne-prijave',
  templateUrl: './krivicne-prijave.component.html',
  styleUrls: ['./krivicne-prijave.component.css']
})
export class KrivicnePrijaveComponent implements OnInit {

  constructor(private router: Router,private tuzilastvoService:TuzilastvoService,public dialog: MatDialog) { }

  krivicnePrijave: KrivicnaPrijava[] = [];

  ngOnInit(): void {
    this.getKrivicnePrijave();
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
}
