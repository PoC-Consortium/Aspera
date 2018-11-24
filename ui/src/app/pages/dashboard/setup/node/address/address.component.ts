import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';

import { BurstNode, constants } from '../../../../../lib/model';
import { NetworkService } from '../../../../../lib/services';
import { NodeService } from '../node.service';
import { SetupService } from '../../setup.service';

@Component({
    selector: 'app-node-setup-address',
    styleUrls: ['./address.component.scss'],
    templateUrl: './address.component.html'
})
export class NodeSetupAddressComponent implements OnInit {
    private nodes: BurstNode[];

    constructor(
        private router: Router,
        private networkService: NetworkService,
        private nodeService: NodeService,
        private setupService: SetupService
    ) { }

    public ngOnInit() {
        this.nodes = constants.nodes;
        for (let node of this.nodes) {
            this.networkService.latency(node).then(ping => {
                node.ping = ping
            })
        }
    }

    public ngOnDestroy() {

    }

    public next() {
        this.nodeService.setNodes(this.nodes.filter(node => node.selected));
        setTimeout(x => {
            this.nodeService.setStepIndex(1)
        }, 0);
    }

}
