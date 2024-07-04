import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { GranicnaPolicijaService } from 'src/app/services/granicna-policija.service';
import { Prelaz } from 'src/app/models/prelaz';
import { Router } from '@angular/router';

@Component({
  selector: 'app-kreiraj-krivicnu-prijavu',
  templateUrl: './kreiraj-krivicnu-prijavu.component.html',
  styleUrls: ['./kreiraj-krivicnu-prijavu.component.css']
})
export class KreirajKrivicnuPrijavuComponent implements OnInit {

  krivicnaPrijavaForm!: FormGroup;
  prelazi: Prelaz[] = [];

  constructor(
    private fb: FormBuilder,
    private granicnaPolicijaService: GranicnaPolicijaService,
    private router: Router,
  ) {}

  ngOnInit(): void {
    this.krivicnaPrijavaForm = this.fb.group({
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
    if (this.krivicnaPrijavaForm.valid) {
      const formValues = this.krivicnaPrijavaForm.value;
      const prelazId = formValues.prelaz;

      const krivicnaPrijava = {
        opis: formValues.opis,
        // Remove prelaz ID from here as it should only send opis in the request body
      };

      this.granicnaPolicijaService.createKrivicnaPrijava(prelazId, krivicnaPrijava).subscribe(
        () => {
          alert('Krivična prijava uspešno kreirana.');
          this.resetForm();
        },
        (error) => {
          console.error(error);
          alert('Došlo je do greške prilikom kreiranja krivične prijave.');
        }
      );
    }
  }

  resetForm(): void {
    this.krivicnaPrijavaForm.reset();
  }

  cancel(): void {
    this.navigateToMainPage();
  }

  navigateToMainPage(): void {
    this.router.navigate(['/Main-Page']);
  }
}
