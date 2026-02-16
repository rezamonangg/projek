import type { ApiResponse } from '../types';

const MOCK_DELAY_MIN = 200;
const MOCK_DELAY_MAX = 800;

function randomDelay(): Promise<void> {
	const delay = Math.random() * (MOCK_DELAY_MAX - MOCK_DELAY_MIN) + MOCK_DELAY_MIN;
	return new Promise((resolve) => setTimeout(resolve, delay));
}

export async function mockResponse<T>(data: T, delay = true): Promise<ApiResponse<T>> {
	if (delay) {
		await randomDelay();
	}
	return {
		data,
		error: null,
		status: 200
	};
}

export async function mockError<T>(message: string, status = 400): Promise<ApiResponse<T>> {
	await randomDelay();
	return {
		data: null,
		error: message,
		status
	};
}

export function generateId(): string {
	return crypto.randomUUID();
}
