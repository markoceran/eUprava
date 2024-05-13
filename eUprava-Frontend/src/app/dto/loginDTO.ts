export class LoginDTO {
    korisnickoIme: string = "";
    lozinka: string = "";

    LoginDTO(korisnickoIme: string, lozinka: string) {
        this.korisnickoIme = korisnickoIme;
        this.lozinka = lozinka;
    }
}