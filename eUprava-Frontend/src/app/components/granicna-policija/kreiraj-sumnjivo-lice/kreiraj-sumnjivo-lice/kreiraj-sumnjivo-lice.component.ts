import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { GranicnaPolicijaService } from 'src/app/services/granicna-policija.service';
import { Prelaz } from 'src/app/models/prelaz';
import { Router } from '@angular/router';
@Component({
  selector: 'app-kreiraj-sumnjivo-lice',
  templateUrl: './kreiraj-sumnjivo-lice.component.html',
  styleUrls: ['./kreiraj-sumnjivo-lice.component.css']
})
export class KreirajSumnjivoLiceComponent implements OnInit {
  sumnjivoLiceForm!: FormGroup;
  prelazi: Prelaz[] = [];

  constructor(
    private fb: FormBuilder,
    private granicnaPolicijaService: GranicnaPolicijaService,
    private router: Router,
  ) { }

  ngOnInit(): void {
    this.sumnjivoLiceForm = this.fb.group({
      opis: ['', Validators.required],
      prelaz: ['', Validators.required],
    });

    this.granicnaPolicijaService.getPrelazi().subscribe(
      (data) => {
        this.prelazi = data;
      },
      (error) => {
        console.error(error);
        alert('Greška pri učitavanju prelaza');
      }
    );
  }

  onSubmit(): void {
    if (this.sumnjivoLiceForm.valid) {
      const formValues = this.sumnjivoLiceForm.value;
      const prelazId = formValues.prelaz;

      const sumnjivoLice = {
        opis: formValues.opis,
        // Remove prelaz ID from here as it should only send opis in the request body
      };

      this.granicnaPolicijaService.createSumnjivoLice(prelazId, sumnjivoLice).subscribe(
        () => {
          alert('Sumnjivo lice uspešno kreirano.');
          this.resetForm();
        },
        (error) => {
          console.error(error);
          alert('Došlo je do greške prilikom kreiranja sumnjivog lica.');
        }
      );
    }
  }

  resetForm(): void {
    this.sumnjivoLiceForm.reset();
  }

  cancel(): void {
    this.navigateToMainPage();
  }

  navigateToMainPage(): void {
    this.router.navigate(['/Main-Page']);
  }

}
