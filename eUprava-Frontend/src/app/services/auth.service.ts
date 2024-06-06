import {HttpClient, HttpHeaders, HttpResponse} from "@angular/common/http";
import { Injectable } from "@angular/core";
import { BehaviorSubject, Observable } from "rxjs";
import { environment } from "src/environments/environment";
import { LoginDTO } from "../dto/loginDTO";
import { jwtDecode } from 'jwt-decode';
import { User } from "../models/user";

@Injectable({
providedIn: 'root'
})
export class AuthService {
  
  private url = "auth";
  constructor(private http: HttpClient) { }

  public Login(loginDTO: LoginDTO): Observable<string> {
    return this.http.post(`${environment.baseApiUrl}/${this.url}/login`, loginDTO, {responseType : 'text'});
  }

  public getUser(userId: any): Observable<User> {
    return this.http.get<User>(`${environment.baseApiUrl}/${this.url}/korisnik/`+ userId);
  }


  isLoggedIn(): boolean {
    if (!localStorage.getItem('authToken')) {
      return false;
    }
    return true;
  }

   private secretKey = 'my_secret_key';

   // Function to parse the JWT token
   parseToken(): any {
    var tokenString = localStorage.getItem("authToken");
    if(tokenString != null){
      try {
       return jwtDecode(tokenString);
     } catch (error) {
       console.error('Error parsing token:', error);
       return null;
     }
    }
    
   }
 
   // Function to extract user type from token
   extractUserType(): string | null {
    const tokenData: any = this.parseToken();
    if (tokenData && tokenData.rola) {
      return tokenData.rola;
    }
    return null;
   }
 
   // Example function to extract claims from token
   extractClaims(): any {
     return this.parseToken();
   }

   getUserIdFromToken(): any {
    const token = localStorage.getItem('authToken');

    if (token) {
      try {
        const payload = token.split('.')[1];
        const decodedPayload = atob(payload);
        const user = JSON.parse(decodedPayload);

        if (user && user.id) {
          const id = user.id;
          return id;
        } else {
          console.error('Invalid user payload:', user);
        }
      } catch (error) {
        console.error('Error decoding token payload:', error);
      }
    } else {
      console.error('Token not found.');
    }

    return null;
  }
 
   // Example function to check if the token is expired
   isTokenExpired(): boolean {
    const tokenData: any = this.parseToken();
    if (tokenData && tokenData.exp) {
      const expiryTimestamp = tokenData.exp * 1000; // Convert to milliseconds
      return Date.now() >= expiryTimestamp;
    }
    return true;
   }

}
