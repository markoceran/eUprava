import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { KrivicnaPrijava } from '../models/krivicnaPrijava';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { ZahtevZaSudskiPostupak } from '../models/zahtevZaSudskiPostupak';
import { ZahtevZaSklapanjeSporazuma } from '../models/zahtevZaSklapanjeSporazuma';
import { Sporazum } from '../models/sporazum';
import { Kanal } from '../models/kanal';
import { Poruka } from '../models/poruka';

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
  public kreirajZahtevZaSudskiPostupak(krivicnaPrijavaId:any, opis:String): Observable<any> {
    return this.http.put<any>(`${environment.baseApiUrl}/${this.url}/kreirajZahtevZaSudskiPostupak/`+ krivicnaPrijavaId, {opis});
  }
  public getZahteviZaSklapanjeSporazumaPoGradjaninu(id:any): Observable<ZahtevZaSklapanjeSporazuma[]> {
    return this.http.get<ZahtevZaSklapanjeSporazuma[]>(`${environment.baseApiUrl}/${this.url}/dobaviZahteveZaSklapanjeSporazumaPoGradjaninu/` + id);
  }
  public prihvatiZahtevZaSklapanjeSporazuma(zahtevId:any): Observable<any> {
    return this.http.put<any>(`${environment.baseApiUrl}/${this.url}/prihvatiZahtevZaSklapanjeSporazuma/`+ zahtevId, {});
  }
  public odbijZahtevZaSklapanjeSporazuma(zahtevId:any): Observable<any> {
    return this.http.put<any>(`${environment.baseApiUrl}/${this.url}/odbijZahtevZaSklapanjeSporazuma/`+ zahtevId, {});
  }
  public kreirajKanal(ime:string, opis:string): Observable<any> {
    return this.http.post<any>(`${environment.baseApiUrl}/${this.url}/kreirajKanal`, {ime,opis});
  }
  public kreirajPoruku(kanalId:any, sadrzaj:string): Observable<any> {
    return this.http.put<any>(`${environment.baseApiUrl}/${this.url}/kreirajPoruku/`+ kanalId, {sadrzaj});
  }
  public getPorukePoKanalu(kanalId:any): Observable<Poruka[]> {
    return this.http.get<Poruka[]>(`${environment.baseApiUrl}/${this.url}/dobaviPorukePoKanalu/` + kanalId);
  }
  public getKanali(): Observable<Kanal[]> {
    return this.http.get<Kanal[]>(`${environment.baseApiUrl}/${this.url}/dobaviKanale`);
  }
}
