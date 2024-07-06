import { Component, OnInit } from '@angular/core';
import { SudService } from '../../services/sud.service';
import { Predmet } from '../../models/predmet';
import { TerminSudjenja } from 'src/app/models/terminSudjenja';
import { Presuda } from 'src/app/models/presuda';

@Component({
  selector: 'app-sud',
  templateUrl: './sud.component.html',
  styleUrls: ['./sud.component.css']
})
export class SudComponent implements OnInit {
  predmeti: Predmet[] = [];
  presude: Presuda[] = [];
  newPredmet: Predmet = {};
  noviTermin: TerminSudjenja = {};
  selectedId: string = "";

  constructor(private sudService: SudService) { }

  ngOnInit(): void {
    this.fetchPresude();
    this.fetchPredmeti();
    this.sudService.createPredmetForZahtev().subscribe(() => {
      this.fetchPredmeti();
    });
  }

  selectPredmet(id:string): void{
    if(this.selectedId == id) {
      this.selectedId = ""
    }
    else {
      this.selectedId = id;
    }
  }

  fetchPredmeti(): void {
    this.sudService.getPredmeti().subscribe((data: Predmet[]) => {
      this.predmeti = data;
    });
  }

  fetchPresude(): void {
    this.sudService.getPresude().subscribe((data: Presuda[]) => {
      this.presude = data;
    });
  }


  createPredmet(): void {
    this.sudService.createPredmet(this.newPredmet).subscribe(() => {
      this.fetchPredmeti();
      this.newPredmet = {};
    });
  }

  zakaziTermin(predmetId: string): void {
    this.sudService.createTermin(this.noviTermin, predmetId).subscribe(() => {
      this.noviTermin = {};
      this.selectedId = ""
    });
  }
}
