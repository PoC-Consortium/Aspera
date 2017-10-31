import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { HttpError } from '../../lib/model';
import { NotificationService } from '../../lib/services';

@Component({
    selector: 'app-register',
    styleUrls: ['./create.component.css'],
    templateUrl: './create.component.html'
})
export class CreateComponent implements OnInit {


    constructor(
        private router: Router,
        private route: ActivatedRoute,
        private notificationService: NotificationService
    ) {

    }

    public ngOnInit() {

    }

    public registerUser() {

    }

    private validateInvite(invite: any) {

    }

    private redirectToLogin() {
        this.router.navigate(['/accounts'], { queryParams: { register: 'success' } });
    }

    private redirectToHome() {
        this.router.navigate(['home']);
    }

    private registrationError(error: HttpError) {
        this.notificationService.error(error.message);
    }

    ngOnDestroy() {

    }

}
