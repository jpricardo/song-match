import z from 'zod';

import { IHttpAdapter } from '@/lib/http';

import {
	GetTrackResponseSchema,
	GetTracksResponseSchema,
	PostAddTrackPayloadSchema,
	PostAddTrackResponseSchema,
	PostFindMatchesResponseSchema,
} from './schema';

export interface ITrackApi {
	getMany(): Promise<z.input<typeof GetTracksResponseSchema>>;
	getOne(id: string): Promise<z.input<typeof GetTrackResponseSchema>>;
	postFindMatches(payload: Uint8Array<ArrayBuffer>): Promise<z.input<typeof PostFindMatchesResponseSchema>>;
	postAddTrack(payload: z.input<typeof PostAddTrackPayloadSchema>): Promise<z.input<typeof PostAddTrackResponseSchema>>;
}

export class TrackApi implements ITrackApi {
	constructor(private readonly baseUrl: string, private readonly httpAdapter: IHttpAdapter) {}

	public async getMany(): Promise<z.input<typeof GetTracksResponseSchema>> {
		return await this.httpAdapter.get(`${this.baseUrl}/tracks`);
	}

	public async getOne(id: string): Promise<z.input<typeof GetTrackResponseSchema>> {
		return await this.httpAdapter.get(`${this.baseUrl}/tracks/${id}`);
	}

	public async postFindMatches(
		payload: Uint8Array<ArrayBuffer>
	): Promise<z.input<typeof PostFindMatchesResponseSchema>> {
		return await this.httpAdapter.post(`${this.baseUrl}/tracks/find`, {
			body: payload,
			headers: { 'Content-Type': 'application/octet-stream' },
		});
	}

	public async postAddTrack(
		payload: z.input<typeof PostAddTrackPayloadSchema>
	): Promise<z.input<typeof PostAddTrackResponseSchema>> {
		return await this.httpAdapter.post(`${this.baseUrl}/tracks`, {
			body: JSON.stringify(payload),
			headers: { 'Content-Type': 'application/json' },
		});
	}
}
