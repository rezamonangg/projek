<script lang="ts">
	import { page } from '$app/stores';
	import { Breadcrumb, Button, Modal, Input } from '$lib/components/common';
	import { tasksApi, boardsApi, membersApi, epicsApi, labelsApi } from '$lib/api';
	import type { Task, Board, User, Epic, Label } from '$lib/types/api';
	import { DndContext, DragOverlay, closestCorners } from '@dnd-kit/core';
	import KanbanColumn from '$lib/components/kanban/KanbanColumn.svelte';
	import TaskCard from '$lib/components/kanban/TaskCard.svelte';

	const projectId = $derived($page.params.id);
	const boardId = $derived($page.params.boardId);

	let board = $state<Board | null>(null);
	let tasks = $state<Task[]>([]);
	let members = $state<User[]>([]);
	let epics = $state<Epic[]>([]);
	let labels = $state<Label[]>([]);
	let loading = $state(true);

	let showCreateModal = $state(false);
	let createTitle = $state('');
	let createDescription = $state('');
	let createStatus = $state<Task['status']>('backlog');
	let createAssigneeId = $state('');
	let createEpicId = $state('');
	let createLoading = $state(false);

	let activeTask = $state<Task | null>(null);

	const columns: { id: Task['status']; label: string }[] = [
		{ id: 'backlog', label: 'Backlog' },
		{ id: 'todo', label: 'To Do' },
		{ id: 'inprogress', label: 'In Progress' },
		{ id: 'done', label: 'Done' }
	];

	async function loadBoard() {
		loading = true;
		const [boardRes, tasksRes, membersRes, epicsRes, labelsRes] = await Promise.all([
			boardsApi.getBoard(boardId),
			tasksApi.listTasks(boardId),
			membersApi.listMembers(),
			epicsApi.listEpics(projectId),
			labelsApi.listLabels('community-1')
		]);

		if (boardRes.data) board = boardRes.data;
		if (tasksRes.data) tasks = tasksRes.data;
		if (membersRes.data) members = membersRes.data.items;
		if (epicsRes.data) epics = epicsRes.data;
		if (labelsRes.data) labels = labelsRes.data;

		loading = false;
	}

	function getTasksByStatus(status: Task['status']) {
		return tasks.filter((t) => t.status === status).sort((a, b) => a.position - b.position);
	}

	function handleDragStart(event: { active: { id: string } }) {
		const task = tasks.find((t) => t.id === event.active.id);
		if (task) {
			activeTask = task;
		}
	}

	async function handleDragEnd(event: { active: { id: string }; over: { id: string } | null }) {
		const { active, over } = event;

		if (!over) {
			activeTask = null;
			return;
		}

		const taskId = active.id as string;
		const task = tasks.find((t) => t.id === taskId);

		if (!task) {
			activeTask = null;
			return;
		}

		const newStatus = over.id as Task['status'];

		if (task.status !== newStatus) {
			await tasksApi.moveTask(taskId, { status: newStatus, position: 0 });
			tasks = tasks.map((t) => (t.id === taskId ? { ...t, status: newStatus } : t));
		}

		activeTask = null;
	}

	async function handleCreate() {
		if (!createTitle.trim()) return;

		createLoading = true;
		const response = await tasksApi.createTask({
			title: createTitle,
			description: createDescription,
			boardId,
			epicId: createEpicId || undefined,
			assigneeId: createAssigneeId || undefined
		});

		createLoading = false;

		if (response.data) {
			if (createStatus !== 'backlog') {
				await tasksApi.moveTask(response.data.id, { status: createStatus, position: 0 });
				response.data.status = createStatus;
			}
			tasks = [...tasks, response.data];
			showCreateModal = false;
			createTitle = '';
			createDescription = '';
			createStatus = 'backlog';
			createAssigneeId = '';
			createEpicId = '';
		}
	}

	$effect(() => {
		loadBoard();
	});
</script>

<svelte:head>
	<title>{board?.name || 'Board'} - Projek</title>
</svelte:head>

<div class="space-y-4">
	<Breadcrumb items={[
		{ label: 'Projects', href: '/projects' },
		{ label: 'Project', href: `/projects/${projectId}` },
		{ label: board?.name || 'Board' }
	]} />

	{#if loading}
		<div class="text-center py-12 text-gray-500">Loading...</div>
	{:else if board}
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-bold text-gray-900">{board.name}</h1>
				{#if board.description}
					<p class="mt-1 text-sm text-gray-500">{board.description}</p>
				{/if}
			</div>
			<Button onclick={() => (showCreateModal = true)}>Add Task</Button>
		</div>

		<DndContext
			collisionDetection={closestCorners}
			ondragstart={handleDragStart}
			ondragend={handleDragEnd}
		>
			<div class="flex gap-4 overflow-x-auto pb-4">
				{#each columns as column}
					<KanbanColumn
						id={column.id}
						title={column.label}
						{tasks}
						{members}
						{epics}
						{labels}
					/>
				{/each}
			</div>

			<DragOverlay>
				{#if activeTask}
					<TaskCard task={activeTask} {members} {epics} {labels} />
				{/if}
			</DragOverlay>
		</DndContext>
	{/if}
</div>

<Modal open={showCreateModal} title="Create Task" onclose={() => (showCreateModal = false)}>
	<div class="space-y-4">
		<Input
			label="Title"
			bind:value={createTitle}
			placeholder="Task title"
			required
		/>

		<div>
			<label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
			<textarea
				bind:value={createDescription}
				rows={3}
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
				placeholder="Task description"
			></textarea>
		</div>

		<div>
			<label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
			<select
				bind:value={createStatus}
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
			>
				{#each columns as column}
					<option value={column.id}>{column.label}</option>
				{/each}
			</select>
		</div>

		<div>
			<label class="block text-sm font-medium text-gray-700 mb-1">Assignee</label>
			<select
				bind:value={createAssigneeId}
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
			>
				<option value="">Unassigned</option>
				{#each members as member}
					<option value={member.id}>{member.firstName} {member.lastName}</option>
				{/each}
			</select>
		</div>

		<div>
			<label class="block text-sm font-medium text-gray-700 mb-1">Epic</label>
			<select
				bind:value={createEpicId}
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
			>
				<option value="">No Epic</option>
				{#each epics as epic}
					<option value={epic.id}>{epic.name}</option>
				{/each}
			</select>
		</div>

		<div class="flex justify-end gap-3">
			<Button variant="secondary" onclick={() => (showCreateModal = false)}>Cancel</Button>
			<Button loading={createLoading} onclick={handleCreate}>Create</Button>
		</div>
	</div>
</Modal>
