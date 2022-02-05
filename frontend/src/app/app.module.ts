import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { MatSliderModule } from '@angular/material/slider';
import {MatToolbarModule } from '@angular/material/toolbar'
import { AppComponent } from './app.component';
import { Routes, RouterModule } from '@angular/router'; 
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ServiceComponent } from './service/service.component';
import { ServiceviewComponent } from './serviceview/serviceview.component';

const routes: Routes = [
  {path:'', component: ServiceviewComponent},
  { path: 'service/:id', component: ServiceComponent }
];

@NgModule({
  declarations: [
    AppComponent,
    ServiceComponent,
    ServiceviewComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    MatSliderModule,
    MatToolbarModule,
    RouterModule.forRoot(routes)
  ],
  providers: [],
  bootstrap: [AppComponent],

  exports: [RouterModule]
  
})
export class AppModule { }