/*
* Copyright 2018 PoC-Consortium
*/

import { Component, OnInit } from "@angular/core";
import { Router } from "@angular/router";
import { CreateService } from "../create.service"

@Component({
    selector: "record",
    moduleId: module.id,
    templateUrl: "./record.component.html",
    styleUrls: ["./record.component.css"]
})
export class RecordComponent implements OnInit {
    private index: number;

    constructor(
        private createService: CreateService,
        private router: Router
    ) {}

    ngOnInit() {
        this.index = 0;
    }

    public onClickNext(e) {
        this.index++;
        if (this.index >= 12) {
            this.index = 0;
            this.createService.setProgress(2)
            this.router.navigate(['create/reproduce'])
        }
    }
}
