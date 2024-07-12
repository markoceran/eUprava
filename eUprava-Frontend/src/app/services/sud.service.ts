import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { Predmet } from '../models/predmet';
import { TerminSudjenja } from '../models/terminSudjenja';
import { Presuda } from '../models/presuda';

@Injectable({
  providedIn: 'root'
})
export class SudService {

  private url = "sud";
  constructor(private http: HttpClient) { }

  public getPredmeti(): Observable<Predmet[]> {
    return this.http.get<Predmet[]>(`${environment.baseApiUrl}/${this.url}/predmeti`);
  }

  getPredmet(id: String): Observable<Predmet> {
    return this.http.get<Predmet>(`${environment.baseApiUrl}/${this.url}/predmeti/${id}`);
  }

  createPredmet(predmet: Predmet) {
    return this.http.post(`${environment.baseApiUrl}/${this.url}/predmeti`, predmet);
  }

  createPredmetForZahtev() {
    return this.http.post(`${environment.baseApiUrl}/${this.url}/predmeti/zahtjevi`, null);
  }

  public getTermini(): Observable<TerminSudjenja[]> {
    return this.http.get<TerminSudjenja[]>(`${environment.baseApiUrl}/${this.url}/termini`);
  }

  getTermin(id: String): Observable<TerminSudjenja> {
    return this.http.get<TerminSudjenja>(`${environment.baseApiUrl}/${this.url}/termini/${id}`);
  }

  createTermin(termin: TerminSudjenja, predmetId: String) {
    return this.http.post(`${environment.baseApiUrl}/${this.url}/termini/${predmetId}`, termin);
  }

  public getPresude(): Observable<Presuda[]> {
    return this.http.get<Presuda[]>(`${environment.baseApiUrl}/${this.url}/presude`);
  }

  getPresuda(id: String): Observable<Presuda> {
    return this.http.get<Presuda>(`${environment.baseApiUrl}/${this.url}/presude/${id}`);
  }

  createPresuda(presuda: Presuda, terminId: String) {
    return this.http.post(`${environment.baseApiUrl}/${this.url}/presude/${terminId}`, presuda);
  }
}
