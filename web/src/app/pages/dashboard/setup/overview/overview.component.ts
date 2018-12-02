import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';

@Component({
    selector: 'app-setup-overview',
    styleUrls: ['./overview.component.css'],
    templateUrl: './overview.component.html'
})
export class SetupOverviewComponent implements OnInit {


    constructor(
        private router: Router,
    ) {}

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

}
