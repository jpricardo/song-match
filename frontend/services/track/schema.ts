export type TrackDTO = { id: number; name: string; matches: number };

export type GetTracksResponse = { tracks: TrackDTO[] };
export type GetTrackResponse = TrackDTO | null;
