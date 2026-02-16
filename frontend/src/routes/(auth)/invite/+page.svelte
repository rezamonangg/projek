<script lang="ts">
	import { page } from '$app/stores';
	import { Button } from '$lib/components/common';
	import { Input } from '$lib/components/common';

	let password = $state('');
	let confirmPassword = $state('');
	let loading = $state(false);
	let errors = $state<{ password?: string; confirmPassword?: string }>({});
	let submitError = $state('');

	const token = $derived($page.url.searchParams.get('token') || '');
	const email = $derived($page.url.searchParams.get('email') || '');

	function validateForm(): boolean {
		errors = {};

		if (!password) {
			errors.password = 'Password is required';
		} else if (password.length < 8) {
			errors.password = 'Password must be at least 8 characters';
		}

		if (!confirmPassword) {
			errors.confirmPassword = 'Please confirm your password';
		} else if (password !== confirmPassword) {
			errors.confirmPassword = 'Passwords do not match';
		}

		return Object.keys(errors).length === 0;
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		submitError = '';

		if (!validateForm()) {
			return;
		}

		loading = true;
		// Will be integrated with mock invitation API
		loading = false;
	}
</script>

<svelte:head>
	<title>Accept Invitation - Projek</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div>
			<h1 class="text-center text-3xl font-bold text-gray-900">Projek</h1>
			<h2 class="mt-6 text-center text-2xl font-semibold text-gray-900">Accept Invitation</h2>
			{#if email}
				<p class="mt-2 text-center text-sm text-gray-600">
					Set your password for <span class="font-medium text-gray-900">{email}</span>
				</p>
			{/if}
		</div>

		{#if !token}
			<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
				Invalid or expired invitation link.
			</div>
		{:else}
			<form class="mt-8 space-y-6" onsubmit={handleSubmit}>
				{#if submitError}
					<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
						{submitError}
					</div>
				{/if}

				<div class="space-y-4">
					<Input
						type="password"
						name="password"
						label="Password"
						bind:value={password}
						error={errors.password}
						required
						placeholder="Create a password"
					/>

					<Input
						type="password"
						name="confirmPassword"
						label="Confirm Password"
						bind:value={confirmPassword}
						error={errors.confirmPassword}
						required
						placeholder="Confirm your password"
					/>
				</div>

				<div>
					<Button type="submit" {loading} class="w-full">
						Accept Invitation
					</Button>
				</div>
			</form>
		{/if}
	</div>
</div>
