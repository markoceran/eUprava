import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ActivatedRoute } from '@angular/router';
import { Poruka } from 'src/app/models/poruka';
import { AuthService } from 'src/app/services/auth.service';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';

@Component({
  selector: 'app-poruke',
  templateUrl: './poruke.component.html',
  styleUrls: ['./poruke.component.css']
})
export class PorukeComponent implements OnInit {

  constructor(private tuzilastvoService:TuzilastvoService, private route: ActivatedRoute,private authService:AuthService,private _snackBar: MatSnackBar) { }

  poruke: Poruka[] = [];
  rolaLogovanogKorisnika: string | null = ""
  inputPolje: string = ""
  kanalId: string = ""

  ngOnInit(): void {
    this.rolaLogovanogKorisnika = this.authService.extractUserType();

    this.route.paramMap.subscribe(params => {
      this.kanalId = params.get('kanalId')!;
      if (this.kanalId != "") {
        this.getPorukePoKanalu(this.kanalId);
      }
    });
  }

  getPorukePoKanalu(kanalId:string): void {
    this.tuzilastvoService.getPorukePoKanalu(kanalId).subscribe(
      (data: Poruka[]) => {
        if(data != null && data.length > 0){
          this.poruke = data;
        }
      },
      (error) => {
        console.error(error);
      }
    );
  }

  posaljiPoruku(){
    if(this.inputPolje != ""){
      this.tuzilastvoService.kreirajPoruku(this.kanalId, this.inputPolje).subscribe(
      (poruka) => {
        this.poruke.push(poruka)
        this.inputPolje = ""
        console.log(poruka)
      },
      (error) => {
        this.openSnackBar(error.message, "");
        console.error(error.message);
        setTimeout(() => {
          window.location.reload();
        }, 2000);
        
      }
    );
    }
  }

  openSnackBar(message: string, action: string) {
    this._snackBar.open(message, action,  {
      duration: 3500
    });
  }

}
