<script lang="ts">
	interface BreadcrumbItem {
		label: string;
		href?: string;
	}

	interface Props {
		items: BreadcrumbItem[];
	}

	let { items }: Props = $props();
</script>

<nav aria-label="Breadcrumb">
	<ol class="flex items-center gap-2 text-sm text-gray-600">
		{#each items as item, index}
			<li class="flex items-center gap-2">
				{#if index > 0}
					<svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
					</svg>
				{/if}
				{#if item.href && index < items.length - 1}
					<a href={item.href} class="hover:text-gray-900 transition-colors">
						{item.label}
					</a>
				{:else}
					<span class="text-gray-900 font-medium">{item.label}</span>
				{/if}
			</li>
		{/each}
	</ol>
</nav>
