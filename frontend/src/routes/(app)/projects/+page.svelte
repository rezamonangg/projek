<script lang="ts">
	import { Breadcrumb, Button, Input, Modal } from '$lib/components/common';
	import { projectsApi } from '$lib/api';
	import type { Project } from '$lib/types/api';

	let projects = $state<Project[]>([]);
	let loading = $state(true);
	let showModal = $state(false);
	let projectName = $state('');
	let projectDescription = $state('');
	let saving = $state(false);
	let errors = $state<{ name?: string }>({});

	async function loadProjects() {
		loading = true;
		const response = await projectsApi.listProjects();
		if (response.data) {
			projects = response.data.items;
		}
		loading = false;
	}

	function validateForm(): boolean {
		errors = {};
		if (!projectName.trim()) {
			errors.name = 'Project name is required';
		}
		return Object.keys(errors).length === 0;
	}

	async function handleCreate() {
		if (!validateForm()) return;

		saving = true;
		const response = await projectsApi.createProject({
			name: projectName,
			description: projectDescription
		});

		saving = false;

		if (response.data) {
			projects = [...projects, response.data];
			showModal = false;
			projectName = '';
			projectDescription = '';
		}
	}

	$effect(() => {
		loadProjects();
	});
</script>

<svelte:head>
	<title>Projects - Projek</title>
</svelte:head>

<div class="space-y-6">
	<Breadcrumb items={[{ label: 'Projects', href: '/projects' }]} />

	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-gray-900">Projects</h1>
			<p class="mt-1 text-sm text-gray-500">Manage your projects</p>
		</div>
		<Button onclick={() => (showModal = true)}>Create Project</Button>
	</div>

	{#if loading}
		<div class="text-center py-12 text-gray-500">Loading...</div>
	{:else if projects.length === 0}
		<div class="text-center py-12 bg-white rounded-lg border border-gray-200">
			<p class="text-gray-500 mb-4">No projects yet.</p>
			<Button onclick={() => (showModal = true)}>Create your first project</Button>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each projects as project}
				<a
					href="/projects/{project.id}"
					class="block bg-white rounded-lg border border-gray-200 p-6 hover:border-blue-300 hover:shadow-sm transition-all"
				>
					<div class="flex items-start justify-between mb-4">
						<h3 class="font-semibold text-gray-900">{project.name}</h3>
						<span class="px-2 py-1 text-xs font-medium rounded-full {project.status === 'active' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'}">
							{project.status}
						</span>
					</div>
					<p class="text-sm text-gray-500 line-clamp-2">{project.description}</p>
					<div class="mt-4 pt-4 border-t border-gray-100 flex items-center justify-between text-xs text-gray-400">
						<span>Created {new Date(project.createdAt).toLocaleDateString()}</span>
					</div>
				</a>
			{/each}
		</div>
	{/if}
</div>

<Modal open={showModal} title="Create Project" onclose={() => (showModal = false)}>
	<div class="space-y-4">
		<Input
			label="Project Name"
			bind:value={projectName}
			error={errors.name}
			placeholder="My Project"
			required
		/>

		<div>
			<label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
			<textarea
				bind:value={projectDescription}
				rows={3}
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
				placeholder="Brief description of the project"
			></textarea>
		</div>

		<div class="flex justify-end gap-3">
			<Button variant="secondary" onclick={() => (showModal = false)}>Cancel</Button>
			<Button loading={saving} onclick={handleCreate}>Create</Button>
		</div>
	</div>
</Modal>
