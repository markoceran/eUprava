import { KrivicnaPrijava } from "./krivicnaPrijava";

export interface ZahtevZaSklapanjeSporazuma {
    id?: string;
    opis?: string;
    uslovi?: string;
    kazna?: string;
    datum?: Date;
    idTuzioca?: string;
    krivicnaPrijava?: KrivicnaPrijava;
    prihvacen?: boolean;
  }