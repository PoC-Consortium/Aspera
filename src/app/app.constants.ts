import { Injectable } from '@angular/core';

@Injectable()
export class Configuration {
    public server: string = 'http://server.url';
    public apiUrl: string = '/api';
    public serverWithApiUrl = this.server + this.apiUrl;
}
