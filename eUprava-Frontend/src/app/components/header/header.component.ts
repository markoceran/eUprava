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
    this.authService.getUser(this.authService.getUserIdFromToken()).subscribe(
      (user: User) => {
        this.role = user.rola;
      },
      (error) => {
        console.error('Error get user data:', error);
      }
    );
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
