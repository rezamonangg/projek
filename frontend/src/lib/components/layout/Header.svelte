<script lang="ts">
	import { authStore } from '$lib/stores/auth';

	interface Props {
		ontoggle: () => void;
	}

	let { ontoggle }: Props = $props();

	let dropdownOpen = $state(false);

	function toggleDropdown() {
		dropdownOpen = !dropdownOpen;
	}

	function closeDropdown() {
		dropdownOpen = false;
	}

	async function handleLogout() {
		closeDropdown();
		await authStore.logout();
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.user-dropdown')) {
			closeDropdown();
		}
	}
</script>

<svelte:window onclick={handleClickOutside} />

<header class="sticky top-0 z-30 flex items-center justify-between h-16 px-4 bg-white border-b border-gray-200 lg:px-6">
	<div class="flex items-center gap-4">
		<button
			type="button"
			class="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg lg:hidden"
			onclick={ontoggle}
			aria-label="Toggle sidebar"
		>
			<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
			</svg>
		</button>

		<nav class="hidden md:flex" aria-label="Breadcrumb">
			<ol class="flex items-center gap-2 text-sm text-gray-600">
				<li><a href="/dashboard" class="hover:text-gray-900">Home</a></li>
			</ol>
		</nav>
	</div>

	<div class="flex items-center gap-4">
		<button
			type="button"
			class="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg"
			aria-label="Notifications"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
			</svg>
		</button>

		<div class="relative user-dropdown">
			<button
				type="button"
				class="flex items-center gap-2 p-2 text-gray-700 hover:bg-gray-100 rounded-lg"
				aria-label="User menu"
				onclick={toggleDropdown}
			>
				<div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center text-white text-sm font-medium">
					{$authStore.user?.firstName?.[0] || 'U'}
				</div>
				<span class="hidden md:block text-sm font-medium">
					{$authStore.user?.firstName || 'User'}
				</span>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
				</svg>
			</button>

			{#if dropdownOpen}
				<div class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-1">
					<div class="px-4 py-2 border-b border-gray-100">
						<p class="text-sm font-medium text-gray-900">{$authStore.user?.firstName} {$authStore.user?.lastName}</p>
						<p class="text-xs text-gray-500">{$authStore.user?.email}</p>
					</div>
					<a
						href="/profile"
						class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
						onclick={closeDropdown}
					>
						Profile
					</a>
					<button
						type="button"
						class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
						onclick={handleLogout}
					>
						Logout
					</button>
				</div>
			{/if}
		</div>
	</div>
</header>
