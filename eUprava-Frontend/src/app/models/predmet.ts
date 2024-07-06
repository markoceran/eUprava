import { ZahtevZaSudskiPostupak } from "./zahtevZaSudskiPostupak";

export interface Predmet{
    id?: string;
    opis?: string;
    datum?: string;
    idSudije?: string;
    zahtev?: ZahtevZaSudskiPostupak;
}