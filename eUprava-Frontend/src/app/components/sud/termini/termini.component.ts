import { Component, OnInit } from '@angular/core';
import { Presuda } from 'src/app/models/presuda';
import { TerminSudjenja } from 'src/app/models/terminSudjenja';
import { SudService } from 'src/app/services/sud.service';

@Component({
  selector: 'app-termini',
  templateUrl: './termini.component.html',
  styleUrls: ['../sud.component.css']
})
export class TerminiComponent implements OnInit {

  termini: TerminSudjenja[] = [];
  novaPresuda: Presuda = {};
  selectedId: string = "";

  constructor(private sudService: SudService) { }

  ngOnInit(): void {
    this.fetchTermini();
  }

  selectTermin(id:string): void{
    if(this.selectedId == id) {
      this.selectedId = ""
    }
    else {
      this.selectedId = id;
    }
  }

  fetchTermini(): void {
    this.sudService.getTermini().subscribe((data: TerminSudjenja[]) => {
      this.termini = data;
    });
  }



  kreirajPresudu(terminId: string): void {
    this.sudService.createPresuda(this.novaPresuda, terminId).subscribe(() => {
      this.novaPresuda = {};
      this.selectedId = ""
    });
  }
}