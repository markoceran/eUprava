import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { Prelaz } from '../models/prelaz';
import { SumnjivoLice } from '../models/sumnjivoLice';
import { KrivicnaPrijava } from '../models/krivicnaPrijava';

@Injectable({
  providedIn: 'root'
})
export class GranicnaPolicijaService {
  private url = 'gp';

  constructor(private http: HttpClient) { }

  // Create Sumnjivo Lice
  public createSumnjivoLice(prelazId: string, sumnjivoLice: any): Observable<any> {
    return this.http.put<any>(`${environment.baseApiUrl}/${this.url}/sumnjivo-lice/new/${prelazId}`, sumnjivoLice);
  }

  // Create Prelaz
  public createPrelaz(prelaz: Prelaz): Observable<any> {
    return this.http.post<any>(`${environment.baseApiUrl}/${this.url}/prelaz/new`, prelaz);
  }

  // Create Krivicna Prijava
  public createKrivicnaPrijava(prelazId: string, krivicnaPrijava: any): Observable<any> {
    return this.http.put<any>(`${environment.baseApiUrl}/${this.url}/krivicna-prijava/new/${prelazId}`, krivicnaPrijava);
  }

  

  // Get Sumnjiva Lica
  public getSumnjivaLica(): Observable<SumnjivoLice[]> {
    return this.http.get<SumnjivoLice[]>(`${environment.baseApiUrl}/${this.url}/sumnjivo-lice/all`);
  }

  // Get Prelazi
  public getPrelazi(): Observable<Prelaz[]> {
    return this.http.get<Prelaz[]>(`${environment.baseApiUrl}/${this.url}/prelaz/all`);
  }

  // Get Krivicne Prijave
  public getKrivicnePrijave(): Observable<KrivicnaPrijava[]> {
    return this.http.get<KrivicnaPrijava[]>(`${environment.baseApiUrl}/${this.url}/krivicna-prijava/all`);
  }
}
