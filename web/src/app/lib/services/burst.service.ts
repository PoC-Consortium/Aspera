import { Injectable } from '@angular/core';
import { GetBlockchainStatusResponse, Block } from '../model';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable()
export class BurstService {

  constructor(private httpService: HttpClient) { }

  public getBlockchainStatus(): Promise<GetBlockchainStatusResponse> {
    return this.httpService.get<GetBlockchainStatusResponse>('getBlockchainStatusResponse').toPromise();
  }

  public getBlock(blockNumber: string): Promise<Block> {
    return this.httpService.get<Block>('getBlock').toPromise();
  }
}
