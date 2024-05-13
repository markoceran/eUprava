import { HttpHeaders } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { LoginDTO } from 'src/app/dto/loginDTO';
import { AuthService } from 'src/app/services/auth.service';
import {MatSnackBar} from "@angular/material/snack-bar";


@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  formGroup: FormGroup = new FormGroup({
    korisnickoIme: new FormControl(''),
    lozinka: new FormControl('')
  });

  constructor(
    private authService: AuthService,
    private router: Router,
    private formBuilder: FormBuilder,
    private _snackBar: MatSnackBar,
  ) { }


  ngOnInit(): void {
    this.formGroup = this.formBuilder.group({
      korisnickoIme: ['', [Validators.required]],
      lozinka: ['', [Validators.required]],
    });
    this.formGroup.setErrors({ unauthenticated: true})
  }


  get loginGroup(): { [key: string]: AbstractControl } {
    return this.formGroup.controls;
  }

  onSubmit() {
      let login: LoginDTO = new LoginDTO();

      login.korisnickoIme = this.formGroup.get('korisnickoIme')?.value;
      login.lozinka = this.formGroup.get('lozinka')?.value;
             
      this.authService.Login(login).subscribe({
        next: (token: string) => {
          localStorage.setItem('authToken', token);
          this.router.navigate(['/Main-Page']);
        },
        error: (error) => {
          this.formGroup.setErrors({ unauthenticated: true });
          this.openSnackBar("Username or password are incorrect!", "");
        }
      });
      
  }

  openSnackBar(message: string, action: string) {
    this._snackBar.open(message, action,  {
      duration: 3500
    });
  }

}
