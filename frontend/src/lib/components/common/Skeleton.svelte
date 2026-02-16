<script lang="ts">
	interface Props {
		width?: string;
		height?: string;
		class?: string;
		rounded?: 'none' | 'sm' | 'md' | 'lg' | 'full';
	}

	let { width = '100%', height = '1rem', class: className = '', rounded = 'md' }: Props = $props();

	const roundedClasses: Record<string, string> = {
		none: 'rounded-none',
		sm: 'rounded-sm',
		md: 'rounded-md',
		lg: 'rounded-lg',
		full: 'rounded-full'
	};
</script>

<div
	class="animate-pulse bg-gray-200 {roundedClasses[rounded]} {className}"
	style="width: {width}; height: {height};"
></div>
