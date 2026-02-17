export function formatDate(dateString: string): string {
	if (!dateString) return '';
	
	let date: Date;
	
	if (dateString.includes('+')) {
		const normalized = dateString.replace(/(\.\d{3})\d+/, '$1');
		date = new Date(normalized);
	} else {
		date = new Date(dateString);
	}
	
	if (isNaN(date.getTime())) {
		return '';
	}
	
	return date.toLocaleDateString();
}

export function formatRelativeTime(dateString: string): string {
	if (!dateString) return '';
	
	let date: Date;
	
	if (dateString.includes('+')) {
		const normalized = dateString.replace(/(\.\d{3})\d+/, '$1');
		date = new Date(normalized);
	} else {
		date = new Date(dateString);
	}
	
	if (isNaN(date.getTime())) {
		return '';
	}
	
	const now = new Date();
	const diffMs = now.getTime() - date.getTime();
	const diffMins = Math.floor(diffMs / 60000);
	const diffHours = Math.floor(diffMs / 3600000);
	const diffDays = Math.floor(diffMs / 86400000);
	
	if (diffMins < 1) return 'just now';
	if (diffMins < 60) return `${diffMins}m ago`;
	if (diffHours < 24) return `${diffHours}h ago`;
	if (diffDays < 7) return `${diffDays}d ago`;
	
	return date.toLocaleDateString();
}
