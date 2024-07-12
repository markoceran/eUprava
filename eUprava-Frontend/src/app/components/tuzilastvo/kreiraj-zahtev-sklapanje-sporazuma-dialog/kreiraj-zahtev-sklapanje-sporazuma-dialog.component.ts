import { Component, Inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';

@Component({
  selector: 'kreiraj-zahtev-sklapanje-sporazuma-dialog',
  templateUrl: './kreiraj-zahtev-sklapanje-sporazuma-dialog.component.html',
  styleUrls: ['./kreiraj-zahtev-sklapanje-sporazuma-dialog.component.css']
})
export class KreirajZahtevSklapanjeSporazumaDialogComponent implements OnInit {
  zahtevForm!: FormGroup;
  krivicnaPrijavaId: any;

  constructor(
    private fb: FormBuilder,
    public dialogRef: MatDialogRef<KreirajZahtevSklapanjeSporazumaDialogComponent>,
    private tuzilastvoService:TuzilastvoService,
    private _snackBar: MatSnackBar,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {this.krivicnaPrijavaId = data.krivicnaPrijavaId;}

  ngOnInit(): void {
    this.zahtevForm = this.fb.group({
      opis: ['', Validators.required],
      uslovi: ['', Validators.required],
      kazna: ['', Validators.required],
    });
  }

  onSubmit(): void {
    console.log('Form submitted');
    if (this.zahtevForm.valid) {
      const formValues = this.zahtevForm.value;
      this.tuzilastvoService.kreirajZahtevZaSklapanjeSporazuma(this.krivicnaPrijavaId, formValues.opis,formValues.uslovi,formValues.kazna).subscribe(
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

