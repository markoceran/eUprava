import {HttpClient, HttpHeaders, HttpResponse} from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { environment } from "src/environments/environment";
import { LoginDTO } from "../dto/loginDTO";
import { jwtDecode } from 'jwt-decode';

@Injectable({
providedIn: 'root'
})
export class AuthService {
  
  private url = "auth";
  private tokenString = localStorage.getItem("authToken");
  constructor(private http: HttpClient) { }

  public Login(loginDTO: LoginDTO): Observable<string> {
    return this.http.post(`${environment.baseApiUrl}/${this.url}/login`, loginDTO, {responseType : 'text'});
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
    if(this.tokenString != null){
      try {
       return jwtDecode(this.tokenString);
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
