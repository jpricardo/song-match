import { ITrackApi } from './api';
import { GetTrackResponse, GetTracksResponse } from './schema';

export interface ITrackService {
	getMany(): Promise<GetTracksResponse>;
	getOne(id: string): Promise<GetTrackResponse>;
}

export class TrackService implements ITrackService {
	constructor(private readonly trackApi: ITrackApi) {}

	public async getMany(): Promise<GetTracksResponse> {
		return await this.trackApi.getMany();
	}

	public async getOne(id: string): Promise<GetTrackResponse> {
		return await this.trackApi.getOne(id);
	}
}
