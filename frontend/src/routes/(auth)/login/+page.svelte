<script lang="ts">
	import { Button } from '$lib/components/common';
	import { Input } from '$lib/components/common';

	let email = $state('');
	let password = $state('');
	let loading = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		loading = true;
		// Will be integrated with mock API in next commit
		loading = false;
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
			<div class="space-y-4">
				<Input
					type="email"
					name="email"
					label="Email address"
					bind:value={email}
					required
					placeholder="you@example.com"
				/>

				<Input
					type="password"
					name="password"
					label="Password"
					bind:value={password}
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
