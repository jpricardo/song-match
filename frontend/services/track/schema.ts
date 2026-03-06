import z from 'zod';

export const FingerprintDTOSchema = z.object({
	timestamp: z.number(),
	peaks: z.number().array(),
});
export type FingerprintDTO = z.output<typeof FingerprintDTOSchema>;

export const TrackDTOSchema = z.object({
	id: z.string(),
	name: z.string(),
	url: z.url(),
	thumbnail: z.url().optional(),
	fingerprints: FingerprintDTOSchema.array(),
	status: z.enum(['', 'failed', 'processing', 'ready']).transform((v) => v || 'ready'),
});
export type TrackDTO = z.output<typeof TrackDTOSchema>;

export const GetTracksDTOSchema = z.object({ tracks: TrackDTOSchema.array() });
export type GetTracksDTO = z.output<typeof GetTracksDTOSchema>;
export const GetTracksResponseSchema = GetTracksDTOSchema;
export type GetTracksResponse = z.output<typeof GetTracksResponseSchema>;

export const GetTrackDTOSchema = TrackDTOSchema.nullable();
export type GetTrackDTO = z.output<typeof GetTrackDTOSchema>;
export const GetTrackResponseSchema = GetTrackDTOSchema;
export type GetTrackResponse = z.output<typeof GetTrackResponseSchema>;

export const PostAddTrackPayloadSchema = TrackDTOSchema.pick({ url: true });
export type PostAddTrackPayload = z.output<typeof PostAddTrackPayloadSchema>;
export const PostAddTrackResponseSchema = TrackDTOSchema;
export type PostAddTrackResponse = z.output<typeof PostAddTrackResponseSchema>;

export const PostFindMatchesPayloadSchema = z.instanceof(Uint8Array<ArrayBuffer>);
export type PostFindMatchesPayload = z.output<typeof PostFindMatchesPayloadSchema>;
export const PostFindMatchesDTOSchema = z.object({ matches: TrackDTOSchema.array() });
export type PostFindMatchesDTO = z.output<typeof PostFindMatchesDTOSchema>;
export const PostFindMatchesResponseSchema = PostFindMatchesDTOSchema;
export type PostFindMatchesResponse = z.output<typeof PostFindMatchesResponseSchema>;
