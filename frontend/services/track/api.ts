import z from 'zod';

import { IHttpAdapter } from '@/lib/http';

import {
	GetTrackResponseSchema,
	GetTracksResponseSchema,
	PostFindMatchesPayloadSchema,
	PostFindMatchesResponseSchema,
} from './schema';

export interface ITrackApi {
	getMany(): Promise<z.input<typeof GetTracksResponseSchema>>;
	getOne(id: string): Promise<z.input<typeof GetTrackResponseSchema>>;
	postFindMatches(
		payload: z.input<typeof PostFindMatchesPayloadSchema>
	): Promise<z.input<typeof PostFindMatchesResponseSchema>>;
}

export class TrackApi implements ITrackApi {
	constructor(private readonly baseUrl: string, private readonly httpAdapter: IHttpAdapter) {}

	public async getMany(): Promise<z.input<typeof GetTracksResponseSchema>> {
		return this.httpAdapter.get(`${this.baseUrl}/tracks`);
	}

	public async getOne(id: string): Promise<z.input<typeof GetTrackResponseSchema>> {
		return this.httpAdapter.get(`${this.baseUrl}/tracks/${id}`);
	}

	public async postFindMatches(
		payload: z.input<typeof PostFindMatchesPayloadSchema>
	): Promise<z.input<typeof PostFindMatchesResponseSchema>> {
		return this.httpAdapter.post(`${this.baseUrl}/tracks/find`, payload);
	}
}
