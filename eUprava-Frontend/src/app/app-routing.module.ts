import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { MainPageComponent } from './components/main-page/main-page.component';
import {LoginGuardService} from "./guards/login-guard.service";
import {RoleGuardService} from "./guards/role-guard.service";
import { KrivicnePrijaveComponent } from './components/tuzilastvo/krivicne-prijave/krivicne-prijave.component';
import { ZahteviSudskiPostupakComponent } from './components/tuzilastvo/zahtevi-sudski-postupak/zahtevi-sudski-postupak.component';
import { ZahteviSklapanjeSporazumaComponent } from './components/tuzilastvo/zahtevi-sklapanje-sporazuma/zahtevi-sklapanje-sporazuma.component';
import { SporazumiComponent } from './components/tuzilastvo/sporazumi/sporazumi.component';
import { KanaliComponent } from './components/tuzilastvo/kanali/kanali.component';
import { PorukeComponent } from './components/tuzilastvo/poruke/poruke.component';
import { KreirajKrivicnuPrijavuComponent } from './components/granicna-policija/kreiraj-krivicnu-prijavu/kreiraj-krivicnu-prijavu/kreiraj-krivicnu-prijavu.component';
import { KreirajPrelazComponent } from './components/granicna-policija/kreiraj-prelaz/kreiraj-prelaz/kreiraj-prelaz.component';


const routes: Routes = [
  {
    path: 'Main-Page',
    component: MainPageComponent
  },
  {
    path: '',
    component: LoginComponent,
    canActivate: [LoginGuardService]
  },
  {
    path: 'krivicnePrijaveTuzilastvo',
    component: KrivicnePrijaveComponent
  },
  {
    path: 'zahteviZaSudskiPostupakTuzilastvo',
    component: ZahteviSudskiPostupakComponent
  },
  {
    path: 'zahteviZaSklapanjeSporazumaTuzilastvo',
    component: ZahteviSklapanjeSporazumaComponent
  },
  {
    path: 'sporazumiTuzilastvo',
    component: SporazumiComponent
  },
  {
    path: 'kanali',
    component: KanaliComponent
  },
  {
    path: 'poruke/:kanalId',
    component: PorukeComponent
  },
  {
    path: 'podnesiKrivicnuPrijavu',
    component: KreirajKrivicnuPrijavuComponent
  },
  {
    path: 'kreirajPrelaz',
    component: KreirajPrelazComponent
  },

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
