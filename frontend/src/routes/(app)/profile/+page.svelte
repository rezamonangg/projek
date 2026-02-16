<script lang="ts">
	import { Breadcrumb, Button, Input } from '$lib/components/common';
	import { membersApi } from '$lib/api';
	import type { User } from '$lib/types/api';

	let user = $state<User | null>(null);
	let firstName = $state('');
	let lastName = $state('');
	let loading = $state(true);
	let saving = $state(false);
	let saveSuccess = $state(false);
	let errors = $state<{ firstName?: string; lastName?: string }>({});

	async function loadProfile() {
		loading = true;
		const response = await membersApi.listMembers();
		if (response.data && response.data.items.length > 0) {
			user = response.data.items[0];
			firstName = user.firstName;
			lastName = user.lastName;
		}
		loading = false;
	}

	function validateForm(): boolean {
		errors = {};
		if (!firstName.trim()) {
			errors.firstName = 'First name is required';
		}
		if (!lastName.trim()) {
			errors.lastName = 'Last name is required';
		}
		return Object.keys(errors).length === 0;
	}

	async function handleSave() {
		if (!user || !validateForm()) return;

		saving = true;
		saveSuccess = false;

		const response = await membersApi.updateProfile(user.id, {
			firstName,
			lastName
		});

		saving = false;

		if (response.data) {
			user = response.data;
			saveSuccess = true;
			setTimeout(() => (saveSuccess = false), 3000);
		}
	}

	$effect(() => {
		loadProfile();
	});
</script>

<svelte:head>
	<title>Profile - Projek</title>
</svelte:head>

<div class="space-y-6">
	<Breadcrumb items={[{ label: 'Profile', href: '/profile' }]} />

	<div>
		<h1 class="text-2xl font-bold text-gray-900">Profile</h1>
		<p class="mt-1 text-sm text-gray-500">Manage your account settings</p>
	</div>

	{#if loading}
		<div class="text-center py-12 text-gray-500">Loading...</div>
	{:else if user}
		<div class="max-w-2xl">
			<div class="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
				{#if saveSuccess}
					<div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg text-sm">
						Profile updated successfully!
					</div>
				{/if}

				<div>
					<h2 class="text-lg font-semibold text-gray-900 mb-4">Personal Information</h2>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<Input
							label="First Name"
							bind:value={firstName}
							error={errors.firstName}
							required
						/>
						<Input
							label="Last Name"
							bind:value={lastName}
							error={errors.lastName}
							required
						/>
					</div>
				</div>

				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
					<input
						type="email"
						value={user.email}
						disabled
						class="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-100 text-gray-500"
					/>
					<p class="mt-1 text-sm text-gray-500">Email cannot be changed</p>
				</div>

				<div class="pt-4 border-t border-gray-200">
					<Button loading={saving} onclick={handleSave}>Save Changes</Button>
				</div>
			</div>

			<div class="mt-6 bg-white rounded-lg border border-gray-200 p-6">
				<h2 class="text-lg font-semibold text-gray-900 mb-4">Change Password</h2>
				<p class="text-sm text-gray-500 mb-4">
					Password change functionality will be available when integrated with the real API.
				</p>
				<Button variant="secondary" disabled>Change Password</Button>
			</div>
		</div>
	{/if}
</div>
