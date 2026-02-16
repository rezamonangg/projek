<script lang="ts">
	import { page } from '$app/stores';
	import { Breadcrumb, Button, Input, Modal } from '$lib/components/common';
	import { wikiApi } from '$lib/api';
	import type { WikiPage } from '$lib/types/api';

	const projectId = $derived($page.params.id);

	let pages = $state<WikiPage[]>([]);
	let loading = $state(true);
	let selectedPage = $state<WikiPage | null>(null);
	let editing = $state(false);
	let showCreateModal = $state(false);
	let createTitle = $state('');
	let createSlug = $state('');
	let createLoading = $state(false);
	let saving = $state(false);

	async function loadPages() {
		loading = true;
		const response = await wikiApi.listWikiPages(projectId);
		if (response.data) {
			pages = response.data;
		}
		loading = false;
	}

	function selectPage(page: WikiPage) {
		selectedPage = page;
		editing = false;
	}

	async function handleCreate() {
		if (!createTitle.trim()) return;

		createLoading = true;
		const slug = createSlug || createTitle.toLowerCase().replace(/[^a-z0-9]+/g, '-');
		const response = await wikiApi.createWikiPage({
			title: createTitle,
			slug,
			content: { type: 'doc', content: [] },
			projectId
		});

		createLoading = false;

		if (response.data) {
			pages = [...pages, response.data];
			showCreateModal = false;
			createTitle = '';
			createSlug = '';
			selectedPage = response.data;
			editing = true;
		}
	}

	async function handleSave() {
		if (!selectedPage) return;
		saving = true;
		const response = await wikiApi.updateWikiPage(selectedPage.id, {
			title: selectedPage.title,
			content: selectedPage.content
		});
		saving = false;
		if (response.data) {
			selectedPage = response.data;
			pages = pages.map((p) => (p.id === response.data!.id ? response.data! : p));
			editing = false;
		}
	}

	function generateSlug(title: string) {
		return title
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, '-')
			.replace(/^-|-$/g, '');
	}

	$effect(() => {
		loadPages();
	});

	$effect(() => {
		if (createTitle && !createSlug) {
			createSlug = generateSlug(createTitle);
		}
	});
</script>

<svelte:head>
	<title>Wiki - Projek</title>
</svelte:head>

<div class="space-y-6">
	<Breadcrumb items={[
		{ label: 'Projects', href: '/projects' },
		{ label: 'Project', href: `/projects/${projectId}` },
		{ label: 'Wiki', href: `/projects/${projectId}/wiki` }
	]} />

	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-gray-900">Wiki</h1>
			<p class="mt-1 text-sm text-gray-500">Project documentation and knowledge base</p>
		</div>
		<Button onclick={() => (showCreateModal = true)}>New Page</Button>
	</div>

	{#if loading}
		<div class="text-center py-12 text-gray-500">Loading...</div>
	{:else}
		<div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
			<div class="bg-white rounded-lg border border-gray-200 p-4">
				<h3 class="text-sm font-semibold text-gray-700 mb-3">Pages</h3>
				{#if pages.length === 0}
					<p class="text-sm text-gray-500">No pages yet.</p>
				{:else}
					<ul class="space-y-1">
						{#each pages as page}
							<li>
								<button
									type="button"
									class="w-full text-left px-3 py-2 rounded-lg text-sm {selectedPage?.id === page.id ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'}"
									onclick={() => selectPage(page)}
								>
									{page.title}
								</button>
							</li>
						{/each}
					</ul>
				{/if}
			</div>

			<div class="lg:col-span-3">
				{#if selectedPage}
					<div class="bg-white rounded-lg border border-gray-200">
						<div class="flex items-center justify-between p-4 border-b border-gray-200">
							{#if editing}
								<input
									type="text"
									bind:value={selectedPage.title}
									class="text-lg font-semibold text-gray-900 border-none focus:outline-none focus:ring-0 flex-1"
								/>
							{:else}
								<h2 class="text-lg font-semibold text-gray-900">{selectedPage.title}</h2>
							{/if}
							<div class="flex gap-2">
								{#if editing}
									<Button size="sm" variant="secondary" onclick={() => (editing = false)}>Cancel</Button>
									<Button size="sm" loading={saving} onclick={handleSave}>Save</Button>
								{:else}
									<Button size="sm" variant="secondary" onclick={() => (editing = true)}>Edit</Button>
								{/if}
							</div>
						</div>
						<div class="p-6">
							{#if editing}
								<textarea
									bind:value={selectedPage.content}
									rows={15}
									class="w-full border-none focus:outline-none focus:ring-0 resize-none"
									placeholder="Start writing..."
								></textarea>
							{:else}
								<div class="prose max-w-none">
									{#if selectedPage.content && Object.keys(selectedPage.content).length > 0}
										{JSON.stringify(selectedPage.content)}
									{:else}
										<p class="text-gray-500">No content yet. Click Edit to start writing.</p>
									{/if}
								</div>
							{/if}
						</div>
					</div>
				{:else}
					<div class="bg-white rounded-lg border border-gray-200 p-12 text-center">
						<p class="text-gray-500">Select a page from the sidebar or create a new one.</p>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>

<Modal open={showCreateModal} title="Create Wiki Page" onclose={() => (showCreateModal = false)}>
	<div class="space-y-4">
		<Input
			label="Title"
			bind:value={createTitle}
			placeholder="Page title"
			required
		/>
		<Input
			label="Slug"
			bind:value={createSlug}
			placeholder="page-url-slug"
		/>
		<div class="flex justify-end gap-3">
			<Button variant="secondary" onclick={() => (showCreateModal = false)}>Cancel</Button>
			<Button loading={createLoading} onclick={handleCreate}>Create</Button>
		</div>
	</div>
</Modal>
