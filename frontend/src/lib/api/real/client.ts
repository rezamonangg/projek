import { PUBLIC_API_URL } from '$env/static/public';
import { ApiError } from '../types';

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

interface RequestOptions {
	method?: HttpMethod;
	body?: unknown;
	headers?: Record<string, string>;
	retries?: number;
}

const MAX_RETRIES = 2;
const RETRY_DELAY = 1000;

async function delay(ms: number): Promise<void> {
	return new Promise((resolve) => setTimeout(resolve, ms));
}

class HttpClient {
	private baseUrl: string;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;
	}

	private async requestWithRetry<T>(
		endpoint: string,
		options: RequestOptions = {},
		retryCount = 0
	): Promise<T> {
		const { method = 'GET', body, headers = {} } = options;

		const config: RequestInit = {
			method,
			headers: {
				'Content-Type': 'application/json',
				...headers
			},
			credentials: 'include'
		};

		if (body && method !== 'GET') {
			config.body = JSON.stringify(body);
		}

		try {
			const response = await fetch(`${this.baseUrl}${endpoint}`, config);

			if (response.status === 401) {
				throw new ApiError(401, 'UNAUTHORIZED', 'Session expired. Please log in again.');
			}

			if (response.status === 403) {
				throw new ApiError(403, 'FORBIDDEN', 'You do not have permission to perform this action.');
			}

			if (response.status === 404) {
				throw new ApiError(404, 'NOT_FOUND', 'The requested resource was not found.');
			}

			if (response.status === 409) {
				throw new ApiError(409, 'CONFLICT', 'A conflict occurred with the current state.');
			}

			if (response.status >= 500 && retryCount < MAX_RETRIES) {
				await delay(RETRY_DELAY * (retryCount + 1));
				return this.requestWithRetry<T>(endpoint, options, retryCount + 1);
			}

			if (!response.ok) {
				let errorMessage = `HTTP ${response.status}`;
				try {
					const errorData = await response.json();
					if (errorData.error?.message) {
						errorMessage = errorData.error.message;
					}
				} catch {
					// Ignore JSON parse errors
				}
				throw new ApiError(response.status, 'API_ERROR', errorMessage);
			}

			if (response.status === 204) {
				return undefined as T;
			}

			const json = await response.json();
			
			if (json && typeof json === 'object' && 'data' in json) {
				return json.data as T;
			}
			
			return json as T;
		} catch (error) {
			if (error instanceof ApiError) {
				throw error;
			}

			if (error instanceof TypeError && error.message.includes('fetch')) {
				if (retryCount < MAX_RETRIES) {
					await delay(RETRY_DELAY * (retryCount + 1));
					return this.requestWithRetry<T>(endpoint, options, retryCount + 1);
				}
				throw new ApiError(0, 'NETWORK_ERROR', 'Unable to connect to the server. Please check your connection.');
			}

			throw new ApiError(0, 'UNKNOWN_ERROR', error instanceof Error ? error.message : 'An unexpected error occurred.');
		}
	}

	async get<T>(endpoint: string): Promise<T> {
		return this.requestWithRetry<T>(endpoint, { method: 'GET' });
	}

	async post<T>(endpoint: string, body?: unknown): Promise<T> {
		return this.requestWithRetry<T>(endpoint, { method: 'POST', body });
	}

	async put<T>(endpoint: string, body?: unknown): Promise<T> {
		return this.requestWithRetry<T>(endpoint, { method: 'PUT', body });
	}

	async patch<T>(endpoint: string, body?: unknown): Promise<T> {
		return this.requestWithRetry<T>(endpoint, { method: 'PATCH', body });
	}

	async delete<T>(endpoint: string): Promise<T> {
		return this.requestWithRetry<T>(endpoint, { method: 'DELETE' });
	}
}

export const httpClient = new HttpClient(PUBLIC_API_URL || '');
