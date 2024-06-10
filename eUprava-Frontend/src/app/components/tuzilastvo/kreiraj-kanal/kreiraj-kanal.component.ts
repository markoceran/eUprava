import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';

@Component({
  selector: 'app-kreiraj-kanal',
  templateUrl: './kreiraj-kanal.component.html',
  styleUrls: ['./kreiraj-kanal.component.css']
})
export class KreirajKanalComponent implements OnInit {

  kanalForm!: FormGroup;

  constructor(
    private fb: FormBuilder,
    public dialogRef: MatDialogRef<KreirajKanalComponent>,
    private tuzilastvoService:TuzilastvoService,
    private _snackBar: MatSnackBar,
  ){}

  ngOnInit(): void {
    this.kanalForm = this.fb.group({
      ime: ['', Validators.required],
      opis: ['', Validators.required],
    });
  }

  onSubmit(): void {
    console.log('Form submitted');
    if (this.kanalForm.valid) {
      const formValues = this.kanalForm.value;
      this.tuzilastvoService.kreirajKanal(formValues.ime,formValues.opis).subscribe(
        (message) => {
          this.dialogRef.close(this.kanalForm.value);
          this.openSnackBar(message.message, "");
          console.log(message.message);
          setTimeout(() => {
            window.location.reload();
          }, 2000);
          
        },
        (error) => {
          this.dialogRef.close(this.kanalForm.value);
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
