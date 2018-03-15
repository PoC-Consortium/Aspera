import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';

import { constants } from '../../../../../lib/model';
import { NodeService } from '../node.service';
import { SetupService } from '../../setup.service';

@Component({
    selector: 'app-node-setup-address',
    styleUrls: ['./address.component.scss'],
    templateUrl: './address.component.html'
})
export class NodeSetupAddressComponent implements OnInit {
    private regions: any;

    constructor(
        private router: Router,
        private setupService: SetupService
    ) { }

    public ngOnInit() {
        this.regions = constants.nodes;
    }

    public ngOnDestroy() {

    }

}
