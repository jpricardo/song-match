import z from 'zod';

function getResponseSchema<T>(shape: z.ZodType<T>) {
	return z.discriminatedUnion('success', [
		z.object({ success: z.literal(false), message: z.string() }),
		z.object({ success: z.literal(true), data: shape }),
	]);
}

export const TrackDTOSchema = z.object({ id: z.number(), name: z.string(), url: z.url() });
export type TrackDTO = z.output<typeof TrackDTOSchema>;

export const GetTracksDTOSchema = z.object({ tracks: TrackDTOSchema.array() });
export type GetTracksDTO = z.output<typeof GetTracksDTOSchema>;
export const GetTracksResponseSchema = getResponseSchema(GetTracksDTOSchema);
export type GetTracksResponse = z.output<typeof GetTracksResponseSchema>;

export const GetTrackDTOSchema = TrackDTOSchema.nullable();
export type GetTrackDTO = z.output<typeof GetTrackDTOSchema>;
export const GetTrackResponseSchema = getResponseSchema(GetTrackDTOSchema);
export type GetTrackResponse = z.output<typeof GetTrackResponseSchema>;

export const PostFindMatchesPayloadSchema = z.object({ content: z.instanceof(Blob) });
export type PostFindMatchesPayload = z.output<typeof PostFindMatchesPayloadSchema>;
export const PostFindMatchesDTOSchema = z.object({ matches: TrackDTOSchema.array() });
export type PostFindMatchesDTO = z.output<typeof PostFindMatchesDTOSchema>;
export const PostFindMatchesResponseSchema = getResponseSchema(PostFindMatchesDTOSchema);
export type PostFindMatchesResponse = z.output<typeof PostFindMatchesResponseSchema>;
