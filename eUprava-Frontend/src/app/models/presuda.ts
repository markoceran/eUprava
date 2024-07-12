import { TerminSudjenja } from "./terminSudjenja";

export interface Presuda{
    id?: string;
    opis?: string;
    datum?: string;
    terminSudjenja?: TerminSudjenja;
    idSudije?: string;
}