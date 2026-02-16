<script lang="ts">
	import { Button } from '$lib/components/common';
	import { Input } from '$lib/components/common';
	import { authStore } from '$lib/stores/auth';

	let email = $state('');
	let password = $state('');
	let loading = $state(false);
	let errors = $state<{ email?: string; password?: string }>({});
	let submitError = $state('');

	function validateForm(): boolean {
		errors = {};

		if (!email) {
			errors.email = 'Email is required';
		} else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
			errors.email = 'Please enter a valid email address';
		}

		if (!password) {
			errors.password = 'Password is required';
		} else if (password.length < 6) {
			errors.password = 'Password must be at least 6 characters';
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
		const success = await authStore.login(email, password);
		loading = false;

		if (!success) {
			submitError = 'Invalid email or password';
		}
	}
</script>

<svelte:head>
	<title>Login - Projek</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div>
			<h1 class="text-center text-3xl font-bold text-gray-900">Projek</h1>
			<h2 class="mt-6 text-center text-2xl font-semibold text-gray-900">Sign in to your account</h2>
			<p class="mt-2 text-center text-sm text-gray-600">
				Or
				<a href="/register" class="font-medium text-blue-600 hover:text-blue-500">
					create a new community
				</a>
			</p>
		</div>

		<form class="mt-8 space-y-6" onsubmit={handleSubmit}>
			{#if submitError}
				<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
					{submitError}
				</div>
			{/if}

			<div class="space-y-4">
				<Input
					type="email"
					name="email"
					label="Email address"
					bind:value={email}
					error={errors.email}
					required
					placeholder="you@example.com"
				/>

				<Input
					type="password"
					name="password"
					label="Password"
					bind:value={password}
					error={errors.password}
					required
					placeholder="Enter your password"
				/>
			</div>

			<div>
				<Button type="submit" {loading} class="w-full">
					Sign in
				</Button>
			</div>
		</form>

		<div class="text-center text-sm text-gray-500">
			<p>Demo credentials: admin@example.com / password123</p>
		</div>
	</div>
</div>
