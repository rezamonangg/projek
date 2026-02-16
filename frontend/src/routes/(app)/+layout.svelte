<script lang="ts">
	import { MainLayout } from '$lib/components/layout';
	import { ToastContainer, setToastContext, ErrorBoundary, OfflineIndicator } from '$lib/components/common';

	let { children } = $props();
</script>

<ErrorBoundary>
	<MainLayout>
		{@render children()}
	</MainLayout>
</ErrorBoundary>

<ToastContainer />
<OfflineIndicator />
