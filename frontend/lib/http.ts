import { URL } from 'url';

// TODO - Support some common options (params, callbacks, etc)
type RequestOptions = Partial<{
	[key: string]: string;
}>;

export interface IHttpAdapter {
	get<T>(url: string | URL, options?: RequestOptions): Promise<T>;
	post<T, K = void>(url: string | URL, payload: T, options?: RequestOptions): Promise<K>;
	put<T, K = void>(url: string | URL, payload: T, options?: RequestOptions): Promise<K>;
	patch<T, K = void>(url: string | URL, payload: T, options?: RequestOptions): Promise<K>;
	delete<T>(url: string | URL, options?: RequestOptions): Promise<T>;
}

export class HttpAdapter implements IHttpAdapter {
	public async get<T>(url: string | URL, options?: RequestOptions): Promise<T> {
		return await fetch(url, { ...options }).then((res) => res.json() as T);
	}

	public async post<T, K = void>(url: string | URL, payload: T, options?: RequestOptions): Promise<K> {
		return await fetch(url, { body: JSON.stringify(payload), ...options }).then((res) => res.json() as K);
	}

	public async put<T, K = void>(url: string | URL, payload: T, options?: RequestOptions): Promise<K> {
		return await fetch(url, { body: JSON.stringify(payload), ...options }).then((res) => res.json() as K);
	}

	public async patch<T, K = void>(url: string | URL, payload: T, options?: RequestOptions): Promise<K> {
		return await fetch(url, { body: JSON.stringify(payload), ...options }).then((res) => res.json() as K);
	}

	public async delete<T>(url: string | URL, options?: RequestOptions): Promise<T> {
		return await fetch(url, { ...options }).then((res) => res.json() as T);
	}
}
