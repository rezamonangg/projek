<script lang="ts">
	import Skeleton from './Skeleton.svelte';

	interface Props {
		rows?: number;
		columns?: number;
	}

	let { rows = 5, columns = 4 }: Props = $props();
</script>

<div class="bg-white rounded-lg border border-gray-200 overflow-hidden">
	<div class="bg-gray-50 px-6 py-3">
		<div class="flex gap-4">
			{#each Array(columns) as _}
				<Skeleton width="100%" height="0.75rem" />
			{/each}
		</div>
	</div>
	<div class="divide-y divide-gray-200">
		{#each Array(rows) as _}
			<div class="px-6 py-4 flex gap-4 items-center">
				{#each Array(columns) as _}
					<Skeleton width="100%" height="1rem" />
				{/each}
			</div>
		{/each}
	</div>
</div>
