import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';

@Component({
    selector: 'app-setup-node',
    styleUrls: ['./node.component.css'],
    templateUrl: './node.component.html'
})
export class NodeComponent implements OnInit {


    constructor(
        private router: Router,
    ) {}

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

}
