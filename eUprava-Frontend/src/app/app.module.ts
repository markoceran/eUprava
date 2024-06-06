import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';

import { MatCardModule } from '@angular/material/card';
import { MatButtonModule} from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatSelectModule} from '@angular/material/select';
import { MatDividerModule } from '@angular/material/divider';
import { MatSnackBarModule } from '@angular/material/snack-bar';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MainPageComponent } from './components/main-page/main-page.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HeaderComponent } from './components/header/header.component';
import { LoginComponent } from './components/login/login.component';
import { AuthInterceptor } from './services/auth.interceptor';
import { NgxCaptchaModule } from 'ngx-captcha';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatInputModule } from '@angular/material/input';
import { MatNativeDateModule } from '@angular/material/core';
import { KrivicnePrijaveComponent } from './components/tuzilastvo/krivicne-prijave/krivicne-prijave.component';
import { ZahteviSudskiPostupakComponent } from './components/tuzilastvo/zahtevi-sudski-postupak/zahtevi-sudski-postupak.component';
import { ZahteviSklapanjeSporazumaComponent } from './components/tuzilastvo/zahtevi-sklapanje-sporazuma/zahtevi-sklapanje-sporazuma.component';
import { SporazumiComponent } from './components/tuzilastvo/sporazumi/sporazumi.component';
import { MatDialogModule } from '@angular/material/dialog';
import { KreirajZahtevSklapanjeSporazumaDialogComponent } from './components/tuzilastvo/kreiraj-zahtev-sklapanje-sporazuma-dialog/kreiraj-zahtev-sklapanje-sporazuma-dialog.component';
import { KreirajZahtevSudskiPostupakComponent } from './components/tuzilastvo/kreiraj-zahtev-sudski-postupak/kreiraj-zahtev-sudski-postupak.component';

@NgModule({
  declarations: [
    AppComponent,
    MainPageComponent,
    HeaderComponent,
    LoginComponent,
    KrivicnePrijaveComponent,
    ZahteviSudskiPostupakComponent,
    ZahteviSklapanjeSporazumaComponent,
    SporazumiComponent,
    KreirajZahtevSklapanjeSporazumaDialogComponent,
    KreirajZahtevSudskiPostupakComponent,
   
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    FormsModule,
    ReactiveFormsModule,
    MatButtonModule,
    MatMenuModule,
    MatToolbarModule,
    MatIconModule,
    MatCardModule,
    MatFormFieldModule,
    MatSelectModule,
    MatDividerModule,
    MatSnackBarModule,
    ReactiveFormsModule,
    NgxCaptchaModule,
    MatDatepickerModule,
    MatInputModule,
    MatNativeDateModule,
    MatDialogModule,
    FormsModule,
    ReactiveFormsModule,
  ],
  providers:
  [{
    provide: HTTP_INTERCEPTORS,
    useClass: AuthInterceptor,
    multi: true,
  }],
  bootstrap: [AppComponent]
})
export class AppModule { }
