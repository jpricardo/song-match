import { HttpAdapter } from '@/lib/http';

import { TrackApi } from './api';
import { TrackService } from './service';

export * from './api';
export * from './schema';
export * from './service';

const httpAdapter = new HttpAdapter();
const trackApi = new TrackApi(process.env.Api_URL!, httpAdapter);
export const trackService = new TrackService(trackApi);
