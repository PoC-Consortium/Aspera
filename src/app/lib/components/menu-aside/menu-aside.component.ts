import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-menu-aside',
  styleUrls: ['./menu-aside.component.css'],
  templateUrl: './menu-aside.component.html'
})
export class MenuAsideComponent implements OnInit {
  private currentUrl: string;

  @Input() private links: Array<any> = [];

  constructor(public router: Router) {
    // getting the current url
    //this.router.events.subscribe((evt) => this.currentUrl = evt.url);
    //this.userService.currentUser.subscribe((user) => this.user = user);
  }

  public ngOnInit() {
    // TODO
  }

}
