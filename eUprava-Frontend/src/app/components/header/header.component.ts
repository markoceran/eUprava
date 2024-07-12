import { Component, OnDestroy, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { User } from 'src/app/models/user';
import { AuthService } from 'src/app/services/auth.service';


@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit{

  role: any;

  constructor(private router: Router,private authService:AuthService) { }

  
  ngOnInit(): void {
    this.role = this.authService.extractUserType();
    console.log(this.role);
  }


  isLoggedIn(): boolean {
    return this.authService.isLoggedIn();
  }


  logout() {
    localStorage.clear();
    this.router.navigate(['']);
  }

}
