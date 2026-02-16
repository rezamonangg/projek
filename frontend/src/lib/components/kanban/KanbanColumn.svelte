<script lang="ts">
	import { useDroppable } from '@dnd-kit/core';
	import type { Task, User, Epic, Label } from '$lib/types/api';
	import TaskCard from './TaskCard.svelte';

	interface Props {
		id: string;
		title: string;
		tasks: Task[];
		members: User[];
		epics: Epic[];
		labels: Label[];
	}

	let { id, title, tasks, members, epics, labels }: Props = $props();

	const { setNodeRef, isOver } = useDroppable({ id });

	const columnTasks = $derived(tasks.filter((t) => t.status === id).sort((a, b) => a.position - b.position));
</script>

<div
	class="flex-shrink-0 w-72 bg-gray-100 rounded-lg p-3"
	class:bg-gray-200={isOver}
	use:setNodeRef
>
	<div class="flex items-center justify-between mb-3">
		<h3 class="text-sm font-semibold text-gray-700">{title}</h3>
		<span class="text-xs text-gray-500 bg-white px-2 py-0.5 rounded-full">{columnTasks.length}</span>
	</div>

	<div class="space-y-2 min-h-[200px]">
		{#each columnTasks as task (task.id)}
			<TaskCard {task} {members} {epics} {labels} />
		{/each}
	</div>
</div>
