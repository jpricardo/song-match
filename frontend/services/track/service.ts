import { ITrackApi } from './api';
import {
	GetTrackDTO,
	GetTrackResponseSchema,
	GetTracksDTO,
	GetTracksResponseSchema,
	PostFindMatchesDTO,
	PostFindMatchesPayload,
	PostFindMatchesPayloadSchema,
	PostFindMatchesResponseSchema,
} from './schema';

export interface ITrackService {
	getMany(): Promise<GetTracksDTO>;
	getOne(id: string): Promise<GetTrackDTO>;
	findMatches(payload: PostFindMatchesPayload): Promise<PostFindMatchesDTO>;
}

export class TrackService implements ITrackService {
	constructor(private readonly trackApi: ITrackApi) {}

	public async getMany(): Promise<GetTracksDTO> {
		const res = await this.trackApi.getMany().then((d) => GetTracksResponseSchema.decode(d));
		if (!res.success) throw new Error(res.message);
		return res.data;
	}

	public async getOne(id: string): Promise<GetTrackDTO> {
		const res = await this.trackApi.getOne(id).then((d) => GetTrackResponseSchema.decode(d));
		if (!res.success) throw new Error(res.message);
		return res.data;
	}

	public async findMatches(payload: PostFindMatchesPayload): Promise<PostFindMatchesDTO> {
		const p = PostFindMatchesPayloadSchema.encode(payload);
		const res = await this.trackApi.postFindMatches(p).then((d) => PostFindMatchesResponseSchema.decode(d));
		if (!res.success) throw new Error(res.message);
		return res.data;
	}
}
