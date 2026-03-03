import { IHttpAdapter } from '@/lib/http';

import { GetTrackResponse, GetTracksResponse } from './schema';

export interface ITrackApi {
	getMany(): Promise<GetTracksResponse>;
	getOne(id: string): Promise<GetTrackResponse>;
}

export class TrackApi implements ITrackApi {
	constructor(private readonly baseUrl: string, private readonly httpAdapter: IHttpAdapter) {}

	public async getMany(): Promise<GetTracksResponse> {
		return this.httpAdapter.get(`${this.baseUrl}/tracks`);
	}

	public async getOne(id: string): Promise<GetTrackResponse> {
		return this.httpAdapter.get(`${this.baseUrl}/tracks/${id}`);
	}
}
