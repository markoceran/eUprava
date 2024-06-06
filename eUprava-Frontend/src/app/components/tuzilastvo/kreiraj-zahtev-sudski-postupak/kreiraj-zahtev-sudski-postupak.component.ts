import { Component, Inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';

@Component({
  selector: 'app-kreiraj-zahtev-sudski-postupak',
  templateUrl: './kreiraj-zahtev-sudski-postupak.component.html',
  styleUrls: ['./kreiraj-zahtev-sudski-postupak.component.css']
})
export class KreirajZahtevSudskiPostupakComponent implements OnInit {

  zahtevForm!: FormGroup;
  krivicnaPrijavaId: any;

  constructor(
    private fb: FormBuilder,
    public dialogRef: MatDialogRef<KreirajZahtevSudskiPostupakComponent>,
    private tuzilastvoService:TuzilastvoService,
    private _snackBar: MatSnackBar,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {this.krivicnaPrijavaId = data.krivicnaPrijavaId;}

  ngOnInit(): void {
    this.zahtevForm = this.fb.group({
      opis: ['', Validators.required]
    });
  }

  onSubmit(): void {
    console.log('Form submitted');
    if (this.zahtevForm.valid) {
      const formValues = this.zahtevForm.value;
      this.tuzilastvoService.kreirajZahtevZaSudskiPostupak(this.krivicnaPrijavaId, formValues.opis).subscribe(
        (message) => {
          this.dialogRef.close(this.zahtevForm.value);
          this.openSnackBar(message.message, "");
          console.log(message.message);
          setTimeout(() => {
            window.location.reload();
          }, 2000);
          
        },
        (error) => {
          this.dialogRef.close(this.zahtevForm.value);
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
