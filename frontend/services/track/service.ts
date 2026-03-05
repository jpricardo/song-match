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
	delete(id: string): Promise<void>;
}

export class TrackService implements ITrackService {
	constructor(private readonly trackApi: ITrackApi) {}

	public async getMany(): Promise<GetTracksDTO> {
		const res = await this.trackApi.getMany();
		if (!res.success) throw new Error(res.message);
		return GetTracksResponseSchema.decode(res.data);
	}

	public async getOne(id: string): Promise<GetTrackDTO> {
		const res = await this.trackApi.getOne(id);
		if (!res.success) throw new Error(res.message);
		return GetTrackResponseSchema.decode(res.data);
	}

	public async findMatches(payload: PostFindMatchesPayload): Promise<PostFindMatchesDTO> {
		const p = PostFindMatchesPayloadSchema.encode(payload);
		const res = await this.trackApi.postFindMatches(p);
		if (!res.success) throw new Error(res.message);
		return PostFindMatchesResponseSchema.decode(res.data);
	}

	public async addTrack(payload: PostAddTrackPayload): Promise<TrackDTO> {
		const p = PostAddTrackPayloadSchema.encode(payload);
		const res = await this.trackApi.postAddTrack(p);
		if (!res.success) throw new Error(res.message);
		return PostAddTrackResponseSchema.decode(res.data);
	}

	public async delete(id: string): Promise<void> {
		const res = await this.trackApi.delete(id);
		if (!res.success) throw new Error(res.message);
		return res.data;
	}
}
