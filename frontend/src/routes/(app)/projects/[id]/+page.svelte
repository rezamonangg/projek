<script lang="ts">
	import { page } from '$app/stores';
	import { Breadcrumb } from '$lib/components/common';
	import { projectsApi, boardsApi, epicsApi } from '$lib/api';
	import type { Project, Board, Epic } from '$lib/types/api';

	let project = $state<Project | null>(null);
	let boards = $state<Board[]>([]);
	let epics = $state<Epic[]>([]);
	let loading = $state(true);
	let activeTab = $state('boards');

	const projectId = $derived($page.params.id);

	async function loadProject() {
		if (!projectId) return;
		loading = true;
		const [projectRes, boardsRes, epicsRes] = await Promise.all([
			projectsApi.getProject(projectId),
			boardsApi.listBoards(projectId),
			epicsApi.listEpics(projectId)
		]);

		if (projectRes.data) {
			project = projectRes.data;
		}
		if (boardsRes.data) {
			boards = boardsRes.data;
		}
		if (epicsRes.data) {
			epics = epicsRes.data;
		}

		loading = false;
	}

	$effect(() => {
		loadProject();
	});

	const tabs = [
		{ id: 'boards', label: 'Boards' },
		{ id: 'epics', label: 'Epics' },
		{ id: 'wiki', label: 'Wiki' },
		{ id: 'settings', label: 'Settings' }
	];
</script>

<svelte:head>
	<title>{project?.name || 'Project'} - Projek</title>
</svelte:head>

<div class="space-y-6">
	<Breadcrumb items={[
		{ label: 'Projects', href: '/projects' },
		{ label: project?.name || 'Project' }
	]} />

	{#if loading}
		<div class="text-center py-12 text-gray-500">Loading...</div>
	{:else if project}
		<div>
			<div class="flex items-start justify-between">
				<div>
					<h1 class="text-2xl font-bold text-gray-900">{project.name}</h1>
					<p class="mt-1 text-sm text-gray-500">{project.description}</p>
				</div>
				<span class="px-3 py-1 text-sm font-medium rounded-full {project.status === 'active' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'}">
					{project.status}
				</span>
			</div>

			<div class="mt-6 border-b border-gray-200">
				<nav class="flex gap-8">
					{#each tabs as tab}
						<button
							type="button"
							class="pb-4 text-sm font-medium border-b-2 transition-colors {activeTab === tab.id
								? 'border-blue-500 text-blue-600'
								: 'border-transparent text-gray-500 hover:text-gray-700'}"
							onclick={() => (activeTab = tab.id)}
						>
							{tab.label}
						</button>
					{/each}
				</nav>
			</div>

			<div class="mt-6">
				{#if activeTab === 'boards'}
					<div class="space-y-4">
						{#if boards.length === 0}
							<div class="text-center py-8 text-gray-500">
								No boards yet. Create your first board!
							</div>
						{:else}
							<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
								{#each boards as board}
									<a
										href="/projects/{projectId}/boards/{board.id}"
										class="block bg-white rounded-lg border border-gray-200 p-4 hover:border-blue-300 transition-colors"
									>
										<h3 class="font-medium text-gray-900">{board.name}</h3>
										{#if board.description}
											<p class="mt-1 text-sm text-gray-500">{board.description}</p>
										{/if}
									</a>
								{/each}
							</div>
						{/if}
					</div>
				{:else if activeTab === 'epics'}
					<div class="space-y-4">
						{#if epics.length === 0}
							<div class="text-center py-8 text-gray-500">
								No epics yet. Create your first epic!
							</div>
						{:else}
							<div class="space-y-3">
								{#each epics as epic}
									<div class="bg-white rounded-lg border border-gray-200 p-4">
										<div class="flex items-center gap-3">
											<div
												class="w-4 h-4 rounded-full"
												style="background-color: {epic.color}"
											></div>
											<h3 class="font-medium text-gray-900">{epic.name}</h3>
										</div>
										{#if epic.description}
											<p class="mt-1 text-sm text-gray-500 ml-7">{epic.description}</p>
										{/if}
									</div>
								{/each}
							</div>
						{/if}
					</div>
				{:else if activeTab === 'wiki'}
					<div class="text-center py-8 text-gray-500">
						Wiki pages will be available here. <a href="/projects/{projectId}/wiki" class="text-blue-600 hover:underline">Go to Wiki</a>
					</div>
				{:else if activeTab === 'settings'}
					<div class="text-center py-8 text-gray-500">
						Project settings will be available here. <a href="/projects/{projectId}/settings" class="text-blue-600 hover:underline">Go to Settings</a>
					</div>
				{/if}
			</div>
		</div>
	{:else}
		<div class="text-center py-12 text-gray-500">Project not found</div>
	{/if}
</div>
