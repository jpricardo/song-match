import z from 'zod';

import { IHttpAdapter } from '@/lib/http';

import {
	GetTrackResponseSchema,
	GetTracksResponseSchema,
	PostAddTrackPayloadSchema,
	PostAddTrackResponseSchema,
	PostFindMatchesResponseSchema,
} from './schema';

type ApiResponse<T> = Promise<{ success: true; data: T } | { success: false; message: string }>;

export interface ITrackApi {
	getMany(): ApiResponse<z.input<typeof GetTracksResponseSchema>>;
	getOne(id: string): ApiResponse<z.input<typeof GetTrackResponseSchema>>;
	postFindMatches(payload: Uint8Array<ArrayBuffer>): ApiResponse<z.input<typeof PostFindMatchesResponseSchema>>;
	postAddTrack(
		payload: z.input<typeof PostAddTrackPayloadSchema>
	): ApiResponse<z.input<typeof PostAddTrackResponseSchema>>;
	delete(id: string): ApiResponse<void>;
}

export class TrackApi implements ITrackApi {
	constructor(private readonly baseUrl: string, private readonly httpAdapter: IHttpAdapter) {}

	public async getMany(): ApiResponse<z.input<typeof GetTracksResponseSchema>> {
		return await this.httpAdapter.get(`${this.baseUrl}/tracks`);
	}

	public async getOne(id: string): ApiResponse<z.input<typeof GetTrackResponseSchema>> {
		return await this.httpAdapter.get(`${this.baseUrl}/tracks/${id}`);
	}

	public async postFindMatches(
		payload: Uint8Array<ArrayBuffer>
	): ApiResponse<z.input<typeof PostFindMatchesResponseSchema>> {
		return await this.httpAdapter.post(`${this.baseUrl}/matches`, {
			body: payload,
			headers: { 'Content-Type': 'application/octet-stream' },
		});
	}

	public async postAddTrack(
		payload: z.input<typeof PostAddTrackPayloadSchema>
	): ApiResponse<z.input<typeof PostAddTrackResponseSchema>> {
		return await this.httpAdapter.post(`${this.baseUrl}/tracks`, {
			body: JSON.stringify(payload),
			headers: { 'Content-Type': 'application/json' },
		});
	}

	public async delete(id: string): ApiResponse<void> {
		return await this.httpAdapter.delete(`${this.baseUrl}/tracks/${id}`);
	}
}
