<script lang="ts">
	import '../app.css';
	import Header from './Header.svelte';
	import Sidebar from './Sidebar.svelte';

	let { children } = $props();

	let sidebarOpen = $state(false);

	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}

	function closeSidebar() {
		sidebarOpen = false;
	}
</script>

<svelte:head>
	<title>Projek</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<Sidebar {sidebarOpen} onclose={closeSidebar} />

	<div class="lg:pl-64">
		<Header ontoggle={toggleSidebar} />

		<main class="p-6">
			{@render children()}
		</main>
	</div>
</div>
