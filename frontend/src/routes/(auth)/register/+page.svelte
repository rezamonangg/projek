<script lang="ts">
	import { Button } from '$lib/components/common';
	import { Input } from '$lib/components/common';

	let communityName = $state('');
	let slug = $state('');
	let adminEmail = $state('');
	let adminPassword = $state('');
	let confirmPassword = $state('');
	let loading = $state(false);
	let errors = $state<Record<string, string>>({});
	let submitError = $state('');

	function generateSlug(name: string): string {
		return name
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, '-')
			.replace(/^-|-$/g, '');
	}

	function validateForm(): boolean {
		errors = {};

		if (!communityName) {
			errors.communityName = 'Community name is required';
		} else if (communityName.length < 3) {
			errors.communityName = 'Community name must be at least 3 characters';
		}

		if (!slug) {
			errors.slug = 'Slug is required';
		} else if (!/^[a-z0-9-]+$/.test(slug)) {
			errors.slug = 'Slug can only contain lowercase letters, numbers, and hyphens';
		}

		if (!adminEmail) {
			errors.adminEmail = 'Email is required';
		} else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(adminEmail)) {
			errors.adminEmail = 'Please enter a valid email address';
		}

		if (!adminPassword) {
			errors.adminPassword = 'Password is required';
		} else if (adminPassword.length < 8) {
			errors.adminPassword = 'Password must be at least 8 characters';
		}

		if (!confirmPassword) {
			errors.confirmPassword = 'Please confirm your password';
		} else if (adminPassword !== confirmPassword) {
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
		// Will be integrated with mock API in next commit
		loading = false;
	}

	$effect(() => {
		if (communityName && !slug) {
			slug = generateSlug(communityName);
		}
	});
</script>

<svelte:head>
	<title>Register Community - Projek</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div>
			<h1 class="text-center text-3xl font-bold text-gray-900">Projek</h1>
			<h2 class="mt-6 text-center text-2xl font-semibold text-gray-900">Create your community</h2>
			<p class="mt-2 text-center text-sm text-gray-600">
				Already have an account?
				<a href="/login" class="font-medium text-blue-600 hover:text-blue-500">
					Sign in
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
					name="communityName"
					label="Community Name"
					bind:value={communityName}
					error={errors.communityName}
					required
					placeholder="Acme Corporation"
				/>

				<Input
					name="slug"
					label="Community Slug"
					bind:value={slug}
					error={errors.slug}
					required
					placeholder="acme-corporation"
				/>

				<div class="border-t border-gray-200 pt-4">
					<p class="text-sm font-medium text-gray-700 mb-4">Admin Account</p>

					<div class="space-y-4">
						<Input
							type="email"
							name="adminEmail"
							label="Admin Email"
							bind:value={adminEmail}
							error={errors.adminEmail}
							required
							placeholder="admin@example.com"
						/>

						<Input
							type="password"
							name="adminPassword"
							label="Password"
							bind:value={adminPassword}
							error={errors.adminPassword}
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
				</div>
			</div>

			<div>
				<Button type="submit" {loading} class="w-full">
					Create Community
				</Button>
			</div>
		</form>
	</div>
</div>
