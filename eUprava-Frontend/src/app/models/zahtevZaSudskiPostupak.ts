import { KrivicnaPrijava } from "./krivicnaPrijava";

export interface ZahtevZaSudskiPostupak {
    id?: string;
    opis?: string;
    datum?: Date;
    idTuzioca?: string;
    krivicnaPrijava?: KrivicnaPrijava;
  }