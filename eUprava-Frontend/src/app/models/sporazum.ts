import { ZahtevZaSklapanjeSporazuma } from "./zahtevZaSklapanjeSporazuma";

export interface Sporazum {
    id?: string;
    zahtev?: ZahtevZaSklapanjeSporazuma;
    datum?: Date;
  }