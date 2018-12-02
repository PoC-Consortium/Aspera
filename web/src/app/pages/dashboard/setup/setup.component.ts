import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';
import { SetupService } from './setup.service';

@Component({
    selector: 'app-setup',
    styleUrls: ['./setup.component.css'],
    templateUrl: './setup.component.html'
})
export class SetupComponent implements OnInit {

    constructor(
        private router: Router,
        private setupService: SetupService
    ) {}

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

}
