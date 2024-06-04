import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Sporazum } from 'src/app/models/sporazum';
import { TuzilastvoService } from 'src/app/services/tuzilastvo.service';

@Component({
  selector: 'app-sporazumi',
  templateUrl: './sporazumi.component.html',
  styleUrls: ['./sporazumi.component.css']
})
export class SporazumiComponent implements OnInit {

  constructor(private router: Router,private tuzilastvoService:TuzilastvoService) { }

  sporazumi: Sporazum[] = [];

  ngOnInit(): void {
    this.getSporazumi();
  }

  getSporazumi(): void {
    this.tuzilastvoService.getSporazumi().subscribe(
      (data) => {
        this.sporazumi = data;
      },
      (error) => {
        console.error(error);
      }
    );
  }
}
