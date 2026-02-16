import { PUBLIC_API_URL } from '$env/static/public';

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

interface RequestOptions {
	method?: HttpMethod;
	body?: unknown;
	headers?: Record<string, string>;
}

interface BackendResponse<T> {
	success: boolean;
	data?: T;
	error?: {
		code: string;
		message: string;
	};
}

class HttpClient {
	private baseUrl: string;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;
	}

	async request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
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

		const response = await fetch(`${this.baseUrl}${endpoint}`, config);

		if (response.status === 401) {
			throw new Error('Unauthorized');
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
			throw new Error(errorMessage);
		}

		if (response.status === 204) {
			return undefined as T;
		}

		const data = await response.json();
		return data;
	}

	async get<T>(endpoint: string): Promise<T> {
		return this.request<T>(endpoint, { method: 'GET' });
	}

	async post<T>(endpoint: string, body?: unknown): Promise<T> {
		return this.request<T>(endpoint, { method: 'POST', body });
	}

	async put<T>(endpoint: string, body?: unknown): Promise<T> {
		return this.request<T>(endpoint, { method: 'PUT', body });
	}

	async patch<T>(endpoint: string, body?: unknown): Promise<T> {
		return this.request<T>(endpoint, { method: 'PATCH', body });
	}

	async delete<T>(endpoint: string): Promise<T> {
		return this.request<T>(endpoint, { method: 'DELETE' });
	}
}

export const httpClient = new HttpClient(PUBLIC_API_URL || '');
