<script lang="ts">
	import { Breadcrumb, Button, Input } from '$lib/components/common';
	import { adminApi, labelsApi } from '$lib/api';
	import type { CommunitySettings, DashboardStats, Label } from '$lib/types/api';

	let settings = $state<CommunitySettings | null>(null);
	let stats = $state<DashboardStats | null>(null);
	let labels = $state<Label[]>([]);
	let loading = $state(true);
	let saving = $state(false);
	let saveSuccess = $state(false);

	let showLabelModal = $state(false);
	let labelName = $state('');
	let labelColor = $state('#3B82F6');
	let labelSaving = $state(false);

	async function loadData() {
		loading = true;
		const [settingsRes, statsRes, labelsRes] = await Promise.all([
			adminApi.getCommunitySettings(),
			adminApi.getDashboardStats(),
			labelsApi.listLabels('community-1')
		]);

		if (settingsRes.data) settings = settingsRes.data;
		if (statsRes.data) stats = statsRes.data;
		if (labelsRes.data) labels = labelsRes.data;

		loading = false;
	}

	async function handleSaveSettings() {
		if (!settings) return;
		saving = true;
		const response = await adminApi.updateCommunitySettings({
			name: settings.name,
			metricsEnabled: settings.metricsEnabled
		});
		saving = false;
		if (response.data) {
			settings = response.data;
			saveSuccess = true;
			setTimeout(() => (saveSuccess = false), 3000);
		}
	}

	async function handleCreateLabel() {
		if (!labelName.trim()) return;
		labelSaving = true;
		const response = await labelsApi.createLabel({
			name: labelName,
			color: labelColor,
			projectId: 'project-1'
		});
		labelSaving = false;
		if (response.data) {
			labels = [...labels, response.data];
			showLabelModal = false;
			labelName = '';
			labelColor = '#3B82F6';
		}
	}

	$effect(() => {
		loadData();
	});
</script>

<svelte:head>
	<title>Admin - Projek</title>
</svelte:head>

<div class="space-y-6">
	<Breadcrumb items={[{ label: 'Admin', href: '/admin' }]} />

	<div>
		<h1 class="text-2xl font-bold text-gray-900">Admin Dashboard</h1>
		<p class="mt-1 text-sm text-gray-500">Manage your community settings</p>
	</div>

	{#if loading}
		<div class="text-center py-12 text-gray-500">Loading...</div>
	{:else}
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<div class="bg-white rounded-lg border border-gray-200 p-6">
				<h2 class="text-lg font-semibold text-gray-900 mb-4">Quick Stats</h2>
				{#if stats}
					<div class="grid grid-cols-2 gap-4">
						<div class="p-4 bg-gray-50 rounded-lg">
							<p class="text-sm text-gray-500">Members</p>
							<p class="text-2xl font-bold text-gray-900">{stats.totalMembers}</p>
						</div>
						<div class="p-4 bg-gray-50 rounded-lg">
							<p class="text-sm text-gray-500">Projects</p>
							<p class="text-2xl font-bold text-gray-900">{stats.totalProjects}</p>
						</div>
						<div class="p-4 bg-gray-50 rounded-lg">
							<p class="text-sm text-gray-500">Tasks</p>
							<p class="text-2xl font-bold text-gray-900">{stats.totalTasks}</p>
						</div>
						<div class="p-4 bg-gray-50 rounded-lg">
							<p class="text-sm text-gray-500">Completed</p>
							<p class="text-2xl font-bold text-gray-900">{stats.tasksByStatus.done}</p>
						</div>
					</div>
				{/if}
			</div>

			<div class="bg-white rounded-lg border border-gray-200 p-6">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-lg font-semibold text-gray-900">Labels</h2>
					<Button size="sm" onclick={() => (showLabelModal = true)}>Add Label</Button>
				</div>
				<div class="flex flex-wrap gap-2">
					{#each labels as label}
						<span
							class="px-3 py-1 text-sm rounded-full"
							style="background-color: {label.color}20; color: {label.color}"
						>
							{label.name}
						</span>
					{/each}
				</div>
			</div>
		</div>

		<div class="bg-white rounded-lg border border-gray-200 p-6">
			<h2 class="text-lg font-semibold text-gray-900 mb-4">Community Settings</h2>
			{#if saveSuccess}
				<div class="mb-4 bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg text-sm">
					Settings saved successfully!
				</div>
			{/if}
			{#if settings}
				<div class="space-y-4 max-w-md">
					<Input
						label="Community Name"
						bind:value={settings.name}
					/>

					<div>
						<label class="block text-sm font-medium text-gray-700 mb-1">Community Slug</label>
						<input
							type="text"
							value={settings.slug}
							disabled
							class="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-100 text-gray-500"
						/>
					</div>

					<div class="flex items-center gap-2">
						<input
							type="checkbox"
							id="metricsEnabled"
							bind:checked={settings.metricsEnabled}
							class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
						/>
						<label for="metricsEnabled" class="text-sm text-gray-700">Enable Metrics Collection</label>
					</div>

					<Button loading={saving} onclick={handleSaveSettings}>Save Settings</Button>
				</div>
			{/if}
		</div>
	{/if}
</div>

{#if showLabelModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" onclick={() => (showLabelModal = false)}>
		<div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6" onclick={(e) => e.stopPropagation()}>
			<h3 class="text-lg font-semibold text-gray-900 mb-4">Create Label</h3>
			<div class="space-y-4">
				<Input
					label="Label Name"
					bind:value={labelName}
					placeholder="e.g., Bug, Feature"
				/>
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Color</label>
					<input
						type="color"
						bind:value={labelColor}
						class="w-full h-10 rounded-lg cursor-pointer"
					/>
				</div>
				<div class="flex justify-end gap-3">
					<Button variant="secondary" onclick={() => (showLabelModal = false)}>Cancel</Button>
					<Button loading={labelSaving} onclick={handleCreateLabel}>Create</Button>
				</div>
			</div>
		</div>
	</div>
{/if}
