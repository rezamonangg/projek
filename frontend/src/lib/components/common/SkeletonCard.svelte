<script lang="ts">
	import Skeleton from './Skeleton.svelte';
</script>

<div class="bg-white rounded-lg border border-gray-200 p-6">
	<div class="space-y-4">
		<div class="flex items-start justify-between">
			<Skeleton width="60%" height="1.25rem" />
			<Skeleton width="4rem" height="1.5rem" rounded="full" />
		</div>
		<div class="space-y-2">
			<Skeleton width="100%" height="0.875rem" />
			<Skeleton width="80%" height="0.875rem" />
		</div>
		<div class="pt-4 border-t border-gray-100 flex justify-between">
			<Skeleton width="6rem" height="0.75rem" />
		</div>
	</div>
</div>
