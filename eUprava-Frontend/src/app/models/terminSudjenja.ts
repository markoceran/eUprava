import { Predmet } from "./predmet";

export interface TerminSudjenja{
    id?: string;
    adresa?: string;
    datum?: string;
    prostorija?: string;
    predmet?: Predmet;
}