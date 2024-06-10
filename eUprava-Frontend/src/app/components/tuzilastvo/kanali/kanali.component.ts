import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { Kanal } from 'src/app/models/kanal';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';
import { KreirajKanalComponent } from '../kreiraj-kanal/kreiraj-kanal.component';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-kanali',
  templateUrl: './kanali.component.html',
  styleUrls: ['./kanali.component.css']
})
export class KanaliComponent implements OnInit {

  constructor(private authService:AuthService,private router: Router, private tuzilastvoService:TuzilastvoService,public dialog: MatDialog) { }

  kanali: Kanal[] = [];
  rolaLogovanogKorisnika: string | null = ""

  ngOnInit(): void {
    this.getKanali();
    this.rolaLogovanogKorisnika = this.authService.extractUserType();
  }

  getKanali(): void {
    this.tuzilastvoService.getKanali().subscribe(
      (data) => {
        this.kanali = data;
      },
      (error) => {
        console.error(error);
      }
    );
  }

  otvoriStranicuZaPoruke(kanalId:string){
    this.router.navigate(['/poruke', kanalId]);
  }

  openKreirajKanalDialog(): void {
    const dialogRef = this.dialog.open(KreirajKanalComponent, {
      width: '400px',
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        console.log('The dialog was closed with result:', result);
        // Optionally handle the result here
      }
    });
  }

}
