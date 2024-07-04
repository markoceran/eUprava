import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Prelaz } from 'src/app/models/prelaz';
import { GranicnaPolicijaService } from 'src/app/services/granicna-policija.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-kreiraj-prelaz',
  templateUrl: './kreiraj-prelaz.component.html',
  styleUrls: ['./kreiraj-prelaz.component.css']
})
export class KreirajPrelazComponent implements OnInit {

  prelazForm!: FormGroup;

  constructor(
    private fb: FormBuilder,
    private granicnaPolicijaService: GranicnaPolicijaService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.prelazForm = this.fb.group({
      imePutnika: ['', Validators.required],
      prezimePutnika: [''],
      JMBGPutnika: [''],
      brojLicneKartePutnika: [''],
      brojPasosaPutnika: [''],
      drzavljanstvoPutnika: [''],
      markaVozila: [''],
      modelVozila: [''],
      svrhaPutovanja: [''],
      odobren: [false]  // Default value for checkbox
    });
  }

  onSubmit(): void {
    if (this.prelazForm.valid) {
      const formValues = this.prelazForm.value;
      const prelaz: Prelaz = {
        imePutnika: formValues.imePutnika,
        prezimePutnika: formValues.prezimePutnika,
        JMBGPutnika: formValues.JMBGPutnika,
        brojLicneKartePutnika: formValues.brojLicneKartePutnika,
        brojPasosaPutnika: formValues.brojPasosaPutnika,
        drzavljanstvoPutnika: formValues.drzavljanstvoPutnika,
        markaVozila: formValues.markaVozila,
        modelVozila: formValues.modelVozila,
        svrhaPutovanja: formValues.svrhaPutovanja,
        odobren: formValues.odobren
      };

      // Call service method to create Prelaz
      this.granicnaPolicijaService.createPrelaz(prelaz).subscribe(
        () => {
          alert('Prelaz uspešno kreiran.');
          this.resetForm();
        },
        (error) => {
          console.error('Došlo je do greške prilikom kreiranja prelaza:', error);
          alert('Došlo je do greške prilikom kreiranja prelaza.');
        }
      );
    }
  }

  resetForm(): void {
    this.prelazForm.reset();
  }

  cancel(): void {
    this.navigateToMainPage();
  }

  navigateToMainPage(): void {
    this.router.navigate(['/Main-Page']); // Change '/main-page' to your actual main page route
  }
}
