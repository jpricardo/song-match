import { URL } from 'url';

type RequestBody = BodyInit | null | undefined;
// TODO - Support some common options (params, callbacks, etc)
type RequestOptions<TBody extends RequestBody = undefined> = {
	signal?: AbortSignal;
	headers?: HeadersInit;
	body?: TBody;
};

export interface IHttpAdapter {
	get<T>(url: string | URL, options?: RequestOptions): Promise<T>;
	post<T extends RequestBody, K = void>(url: string | URL, options?: RequestOptions<T>): Promise<K>;
	put<T extends RequestBody, K = void>(url: string | URL, options?: RequestOptions<T>): Promise<K>;
	patch<T extends RequestBody, K = void>(url: string | URL, options?: RequestOptions<T>): Promise<K>;
	delete<T>(url: string | URL, options?: RequestOptions): Promise<T>;
}

export class HttpAdapter implements IHttpAdapter {
	public async get<T>(url: string | URL, options?: RequestOptions): Promise<T> {
		return await fetch(url, { method: 'GET', ...options }).then((res) => res.json() as T);
	}

	public async post<T extends RequestBody, K = void>(url: string | URL, options?: RequestOptions<T>): Promise<K> {
		return await fetch(url, { method: 'POST', ...options }).then((res) => res.json() as K);
	}

	public async put<T extends RequestBody, K = void>(url: string | URL, options?: RequestOptions<T>): Promise<K> {
		return await fetch(url, { method: 'PUT', ...options }).then((res) => res.json() as K);
	}

	public async patch<T extends RequestBody, K = void>(url: string | URL, options?: RequestOptions<T>): Promise<K> {
		return await fetch(url, { method: 'PATCH', ...options }).then((res) => res.json() as K);
	}

	public async delete<T>(url: string | URL, options?: RequestOptions): Promise<T> {
		return await fetch(url, { method: 'DELETE', ...options }).then((res) => res.json() as T);
	}
}
