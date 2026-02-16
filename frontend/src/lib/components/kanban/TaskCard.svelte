<script lang="ts">
	import { useDraggable } from '@dnd-kit/core';
	import type { Task, User, Epic, Label } from '$lib/types/api';

	interface Props {
		task: Task;
		members: User[];
		epics: Epic[];
		labels: Label[];
	}

	let { task, members, epics, labels }: Props = $props();

	const { setNodeRef, attributes, listeners, isDragging } = useDraggable({
		id: task.id
	});

	const assignee = $derived(members.find((m) => m.id === task.assigneeId));
	const epic = $derived(epics.find((e) => e.id === task.epicId));
	const taskLabels = $derived(labels.filter((l) => task.labels.includes(l.id)));
</script>

<div
	use:setNodeRef
	{...attributes}
	{...listeners}
	class="bg-white rounded-lg border border-gray-200 p-3 shadow-sm cursor-grab hover:border-blue-300 transition-colors"
	class:hidden={isDragging}
	class:opacity-50={isDragging}
>
	<div class="flex items-start justify-between gap-2">
		<h4 class="text-sm font-medium text-gray-900">{task.title}</h4>
	</div>

	{#if task.description}
		<p class="mt-1 text-xs text-gray-500 line-clamp-2">{task.description}</p>
	{/if}

	{#if taskLabels.length > 0}
		<div class="mt-2 flex flex-wrap gap-1">
			{#each taskLabels.slice(0, 3) as label}
				<span
					class="px-2 py-0.5 text-xs rounded-full"
					style="background-color: {label.color}20; color: {label.color}"
				>
					{label.name}
				</span>
			{/each}
			{#if taskLabels.length > 3}
				<span class="px-2 py-0.5 text-xs rounded-full bg-gray-100 text-gray-600">
					+{taskLabels.length - 3}
				</span>
			{/if}
		</div>
	{/if}

	<div class="mt-3 flex items-center justify-between">
		{#if epic}
			<div class="flex items-center gap-1">
				<div class="w-2 h-2 rounded-full" style="background-color: {epic.color}"></div>
				<span class="text-xs text-gray-500">{epic.name}</span>
			</div>
		{:else}
			<div></div>
		{/if}

		{#if assignee}
			<div
				class="w-6 h-6 bg-blue-600 rounded-full flex items-center justify-center text-white text-xs font-medium"
				title="{assignee.firstName} {assignee.lastName}"
			>
				{assignee.firstName[0]}{assignee.lastName[0]}
			</div>
		{/if}
	</div>
</div>
