import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';

@Component({
    selector: 'app-account-create',
    styleUrls: ['./create.component.scss'],
    templateUrl: './create.component.html'
})
export class AccountCreateComponent implements OnInit {


    constructor(
        private router: Router,
    ) {}

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

}
