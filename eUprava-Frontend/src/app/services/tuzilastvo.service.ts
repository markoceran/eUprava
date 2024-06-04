import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { KrivicnaPrijava } from '../models/krivicnaPrijava';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { ZahtevZaSudskiPostupak } from '../models/zahtevZaSudskiPostupak';
import { ZahtevZaSklapanjeSporazuma } from '../models/zahtevZaSklapanjeSporazuma';
import { Sporazum } from '../models/sporazum';

@Injectable({
  providedIn: 'root'
})
export class TuzilastvoService {

  private url = "tuzilastvo";
  constructor(private http: HttpClient) { }

  public getKrivicnePrijave(): Observable<KrivicnaPrijava[]> {
    return this.http.get<KrivicnaPrijava[]>(`${environment.baseApiUrl}/${this.url}/krivicnePrijave`);
  }
  public getZahteviZaSudskiPostupak(): Observable<ZahtevZaSudskiPostupak[]> {
    return this.http.get<ZahtevZaSudskiPostupak[]>(`${environment.baseApiUrl}/${this.url}/dobaviZahteveZaSudskiPostupak`);
  }
  public getZahteviZaSklapanjeSporazuma(): Observable<ZahtevZaSklapanjeSporazuma[]> {
    return this.http.get<ZahtevZaSklapanjeSporazuma[]>(`${environment.baseApiUrl}/${this.url}/dobaviZahteveZaSklapanjeSporazuma`);
  }
  public getSporazumi(): Observable<Sporazum[]> {
    return this.http.get<Sporazum[]>(`${environment.baseApiUrl}/${this.url}/dobaviSporazume`);
  }
  public kreirajZahtevZaSklapanjeSporazuma(krivicnaPrijavaId:any, opis:String, uslovi:String, kazna:String): Observable<any> {
    return this.http.put<any>(`${environment.baseApiUrl}/${this.url}/kreirajZahtevZaSklapanjeSporazuma/`+ krivicnaPrijavaId, {opis,uslovi,kazna});
  }
}
