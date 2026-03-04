import { ITrackApi } from './api';
import {
	GetTrackDTO,
	GetTrackResponseSchema,
	GetTracksDTO,
	GetTracksResponseSchema,
	PostAddTrackPayload,
	PostAddTrackPayloadSchema,
	PostAddTrackResponseSchema,
	PostFindMatchesDTO,
	PostFindMatchesPayload,
	PostFindMatchesPayloadSchema,
	PostFindMatchesResponseSchema,
	TrackDTO,
} from './schema';

export interface ITrackService {
	getMany(): Promise<GetTracksDTO>;
	getOne(id: string): Promise<GetTrackDTO>;
	findMatches(payload: PostFindMatchesPayload): Promise<PostFindMatchesDTO>;
	addTrack(payload: PostAddTrackPayload): Promise<TrackDTO>;
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

	public async addTrack(payload: PostAddTrackPayload): Promise<TrackDTO> {
		const p = PostAddTrackPayloadSchema.encode(payload);
		const res = await this.trackApi.postAddTrack(p).then((d) => PostAddTrackResponseSchema.decode(d));
		if (!res.success) throw new Error(res.message);
		return res.data;
	}
}
