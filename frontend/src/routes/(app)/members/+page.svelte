<script lang="ts">
	import { Breadcrumb, Button, Input, Modal, SkeletonTable } from '$lib/components/common';
	import { membersApi } from '$lib/api';
	import type { User } from '$lib/types/api';

	let members = $state<User[]>([]);
	let loading = $state(true);
	let searchQuery = $state('');
	let showInviteModal = $state(false);
	let inviteEmail = $state('');
	let inviteRole = $state<'admin' | 'member'>('member');
	let inviteLoading = $state(false);
	let inviteError = $state('');
	let inviteSuccess = $state(false);

	async function loadMembers() {
		loading = true;
		const response = await membersApi.listMembers();
		if (response.data) {
			members = response.data.items;
		}
		loading = false;
	}

	async function handleInvite() {
		if (!inviteEmail) return;

		inviteLoading = true;
		inviteError = '';

		const response = await membersApi.inviteMember({ email: inviteEmail, role: inviteRole });

		inviteLoading = false;

		if (response.error) {
			inviteError = response.error;
		} else {
			inviteSuccess = true;
			inviteEmail = '';
			showInviteModal = false;
			loadMembers();
		}
	}

	$effect(() => {
		loadMembers();
	});

	const filteredMembers = $derived(
		members.filter(
			(m) =>
				m.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
				m.firstName.toLowerCase().includes(searchQuery.toLowerCase()) ||
				m.lastName.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);
</script>

<svelte:head>
	<title>Members - Projek</title>
</svelte:head>

<div class="space-y-6">
	<Breadcrumb items={[{ label: 'Members', href: '/members' }]} />

	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-gray-900">Members</h1>
			<p class="mt-1 text-sm text-gray-500">Manage your community members</p>
		</div>
		<Button onclick={() => (showInviteModal = true)}>Invite Member</Button>
	</div>

	<div class="bg-white rounded-lg border border-gray-200">
		<div class="p-4 border-b border-gray-200">
			<Input
				type="search"
				placeholder="Search members..."
				bind:value={searchQuery}
			/>
		</div>

		{#if loading}
			<div class="p-4"><SkeletonTable rows={5} columns={3} /></div>
		{:else if filteredMembers.length === 0}
			<div class="p-12 text-center text-gray-500">
				{searchQuery ? 'No members found matching your search.' : 'No members yet.'}
			</div>
		{:else}
			<table class="w-full">
				<thead class="bg-gray-50">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Member</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Joined</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200">
					{#each filteredMembers as member}
						<tr class="hover:bg-gray-50">
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="flex items-center gap-3">
									<div class="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center text-white font-medium">
										{member.firstName[0]}{member.lastName[0]}
									</div>
									<div>
										<div class="font-medium text-gray-900">{member.firstName} {member.lastName}</div>
										<div class="text-sm text-gray-500">{member.email}</div>
									</div>
								</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span class="px-2 py-1 text-xs font-medium rounded-full {member.role === 'admin' ? 'bg-purple-100 text-purple-700' : 'bg-gray-100 text-gray-700'}">
									{member.role}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
								{new Date(member.createdAt).toLocaleDateString()}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{/if}
	</div>
</div>

<Modal open={showInviteModal} title="Invite Member" onclose={() => (showInviteModal = false)}>
	{#if inviteError}
		<div class="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
			{inviteError}
		</div>
	{/if}

	<div class="space-y-4">
		<Input
			type="email"
			label="Email Address"
			bind:value={inviteEmail}
			placeholder="member@example.com"
			required
		/>

		<div>
			<label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
			<select
				bind:value={inviteRole}
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
			>
				<option value="member">Member</option>
				<option value="admin">Admin</option>
			</select>
		</div>

		<div class="flex justify-end gap-3">
			<Button variant="secondary" onclick={() => (showInviteModal = false)}>Cancel</Button>
			<Button loading={inviteLoading} onclick={handleInvite}>Send Invite</Button>
		</div>
	</div>
</Modal>

{#if inviteSuccess}
	<div class="fixed bottom-4 right-4 bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg">
		Invitation sent successfully!
	</div>
{/if}
