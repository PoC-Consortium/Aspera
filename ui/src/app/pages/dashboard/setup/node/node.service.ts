import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { BurstNode } from '../../../../lib/model';

@Injectable()
export class NodeService {
    private nodes: BurstNode[];
    private stepIndex: number;

    constructor() {
        this.reset();
    }

    public setNodes(nodes: BurstNode[]) {
        this.nodes = nodes;
    }

    public getNodes() {
        return this.nodes;
    }

    public setStepIndex(index: number) {
        this.stepIndex = index;
    }

    public getStepIndex() : number {
        return this.stepIndex;
    }

    public reset() {
        this.stepIndex = 0;
        this.nodes = [];
    }
}
