import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpError } from '../../lib/model';
import { ActivatedRoute, Router, Params } from '@angular/router';

@Component({
    selector: 'particle-accounts',
    styleUrls: ['./accounts.component.css'],
    templateUrl: './accounts.component.html'
})
export class AccountsComponent implements OnInit {
    public email: string;
    public password: string;

    public err: boolean;
    public errMsg: string;

    public showRegister: boolean = false;

    constructor(
        public router: Router,
        private activatedRoute: ActivatedRoute
    ) {
    }

    public ngOnInit() {
        window.dispatchEvent(new Event('resize'));

        this.activatedRoute.queryParams.subscribe((params: Params) => {
            this.showRegister = params['register'] == "success";
        });
    }

    public login() {

    }

    public setLoggedInUser(token: string) {

    }

    public redirect() {
        this.router.navigate(['/dashboard/home']);
    }

    public loginError(error: HttpError) {
        this.err = true;
        this.errMsg = "Wrong email or Password!"
    }

    public noError() {
        this.err = false;
        this.errMsg = '';
    }


}
