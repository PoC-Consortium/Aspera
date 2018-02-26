import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';

@Component({
    selector: 'app-account-create-seed',
    styleUrls: ['./seed.component.scss'],
    templateUrl: './seed.component.html'
})
export class AccountCreateSeedComponent implements OnInit {
    private seedLimit: number = 1024;
    private seed: any[] = [];
    private update: boolean = false;
    private progress = 0;

    constructor(
        private router: Router,
    ) {}

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

    public movement(e) {
        this.seed.push([e.clientX, e.clientY, new Date()]);
        if (!this.update) {
            this.update = true
            setTimeout(() => {
                this.progress = this.seed.length / this.seedLimit * 100;
                this.update = false;
            }, 100)
        }
    }
}
