export interface ApiResponse<T> {
	data: T | null;
	error: string | null;
	status: number;
}

export interface PaginatedResponse<T> {
	items: T[];
	total: number;
	page: number;
	pageSize: number;
	totalPages: number;
}

export interface ApiClient {
	get<T>(url: string): Promise<ApiResponse<T>>;
	post<T>(url: string, body?: unknown): Promise<ApiResponse<T>>;
	put<T>(url: string, body?: unknown): Promise<ApiResponse<T>>;
	delete<T>(url: string): Promise<ApiResponse<T>>;
}
