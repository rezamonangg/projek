<script lang="ts">
	import { MainLayout } from '$lib/components/layout';
	import { ToastContainer, setToastContext } from '$lib/components/common';

	let { children } = $props();
</script>

<MainLayout>
	{@render children()}
</MainLayout>

<ToastContainer />
