import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';
import { NodeService } from './node.service';

@Component({
    selector: 'app-node-setup',
    styleUrls: ['./node.component.css'],
    templateUrl: './node.component.html'
})
export class NodeSetupComponent implements OnInit {


    constructor(
        private router: Router,
        private nodeService: NodeService
    ) {}

    public ngOnInit() {

    }

    public ngOnDestroy() {

    }

}
