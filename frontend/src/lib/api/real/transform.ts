export function toCamelCase<T>(obj: unknown): T {
	if (obj === null || obj === undefined) {
		return obj as T;
	}

	if (Array.isArray(obj)) {
		return obj.map((item) => toCamelCase<unknown>(item)) as T;
	}

	if (typeof obj !== 'object') {
		return obj as T;
	}

	const result: Record<string, unknown> = {};
	for (const [key, value] of Object.entries(obj as Record<string, unknown>)) {
		const camelKey = key.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
		result[camelKey] = toCamelCase(value);
	}
	return result as T;
}

export function toSnakeCase<T>(obj: unknown): T {
	if (obj === null || obj === undefined) {
		return obj as T;
	}

	if (Array.isArray(obj)) {
		return obj.map((item) => toSnakeCase<unknown>(item)) as T;
	}

	if (typeof obj !== 'object') {
		return obj as T;
	}

	const result: Record<string, unknown> = {};
	for (const [key, value] of Object.entries(obj as Record<string, unknown>)) {
		const snakeKey = key.replace(/[A-Z]/g, (letter) => `_${letter.toLowerCase()}`);
		result[snakeKey] = toSnakeCase(value);
	}
	return result as T;
}
