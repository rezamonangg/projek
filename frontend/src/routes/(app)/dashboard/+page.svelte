<script lang="ts">
	import { Breadcrumb } from '$lib/components/common';
	import { adminApi, projectsApi } from '$lib/api';

	let stats = $state({
		totalMembers: 0,
		totalProjects: 0,
		totalTasks: 0,
		tasksByStatus: { backlog: 0, todo: 0, inprogress: 0, done: 0 }
	});

	let projects = $state<{ id: string; name: string; description: string; status: string }[]>([]);
	let loading = $state(true);

	async function loadDashboard() {
		loading = true;
		const [statsRes, projectsRes] = await Promise.all([
			adminApi.getDashboardStats(),
			projectsApi.listProjects()
		]);

		if (statsRes.data) {
			stats = statsRes.data;
		}

		if (projectsRes.data) {
			projects = projectsRes.data.items.slice(0, 5);
		}

		loading = false;
	}

	$effect(() => {
		loadDashboard();
	});
</script>

<svelte:head>
	<title>Dashboard - Projek</title>
</svelte:head>

<div class="space-y-6">
	<Breadcrumb items={[{ label: 'Dashboard', href: '/dashboard' }]} />

	<div>
		<h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
		<p class="mt-1 text-sm text-gray-500">Welcome back! Here's an overview of your workspace.</p>
	</div>

	{#if loading}
		<div class="text-center py-12 text-gray-500">Loading...</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
			<div class="bg-white rounded-lg border border-gray-200 p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-500">Total Members</p>
						<p class="text-2xl font-bold text-gray-900">{stats.totalMembers}</p>
					</div>
					<div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
						<svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
						</svg>
					</div>
				</div>
			</div>

			<div class="bg-white rounded-lg border border-gray-200 p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-500">Total Projects</p>
						<p class="text-2xl font-bold text-gray-900">{stats.totalProjects}</p>
					</div>
					<div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
						<svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
						</svg>
					</div>
				</div>
			</div>

			<div class="bg-white rounded-lg border border-gray-200 p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-500">Total Tasks</p>
						<p class="text-2xl font-bold text-gray-900">{stats.totalTasks}</p>
					</div>
					<div class="w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center">
						<svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
						</svg>
					</div>
				</div>
			</div>

			<div class="bg-white rounded-lg border border-gray-200 p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-500">Completed</p>
						<p class="text-2xl font-bold text-gray-900">{stats.tasksByStatus.done}</p>
					</div>
					<div class="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
						<svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
				</div>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<div class="bg-white rounded-lg border border-gray-200 p-6">
				<h2 class="text-lg font-semibold text-gray-900 mb-4">Recent Projects</h2>
				{#if projects.length === 0}
					<p class="text-gray-500 text-sm">No projects yet. Create your first project!</p>
				{:else}
					<div class="space-y-3">
						{#each projects as project}
							<a
								href="/projects/{project.id}"
								class="block p-4 rounded-lg border border-gray-200 hover:border-blue-300 hover:bg-blue-50 transition-colors"
							>
								<div class="flex items-center justify-between">
									<div>
										<h3 class="font-medium text-gray-900">{project.name}</h3>
										<p class="text-sm text-gray-500">{project.description}</p>
									</div>
									<span class="px-2 py-1 text-xs font-medium rounded-full {project.status === 'active' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'}">
										{project.status}
									</span>
								</div>
							</a>
						{/each}
					</div>
				{/if}
			</div>

			<div class="bg-white rounded-lg border border-gray-200 p-6">
				<h2 class="text-lg font-semibold text-gray-900 mb-4">Task Distribution</h2>
				<div class="space-y-4">
					<div>
						<div class="flex justify-between text-sm mb-1">
							<span class="text-gray-600">Backlog</span>
							<span class="font-medium">{stats.tasksByStatus.backlog}</span>
						</div>
						<div class="w-full bg-gray-200 rounded-full h-2">
							<div
								class="bg-gray-500 h-2 rounded-full"
								style="width: {(stats.tasksByStatus.backlog / stats.totalTasks) * 100}%"
							></div>
						</div>
					</div>
					<div>
						<div class="flex justify-between text-sm mb-1">
							<span class="text-gray-600">To Do</span>
							<span class="font-medium">{stats.tasksByStatus.todo}</span>
						</div>
						<div class="w-full bg-gray-200 rounded-full h-2">
							<div
								class="bg-blue-500 h-2 rounded-full"
								style="width: {(stats.tasksByStatus.todo / stats.totalTasks) * 100}%"
							></div>
						</div>
					</div>
					<div>
						<div class="flex justify-between text-sm mb-1">
							<span class="text-gray-600">In Progress</span>
							<span class="font-medium">{stats.tasksByStatus.inprogress}</span>
						</div>
						<div class="w-full bg-gray-200 rounded-full h-2">
							<div
								class="bg-yellow-500 h-2 rounded-full"
								style="width: {(stats.tasksByStatus.inprogress / stats.totalTasks) * 100}%"
							></div>
						</div>
					</div>
					<div>
						<div class="flex justify-between text-sm mb-1">
							<span class="text-gray-600">Done</span>
							<span class="font-medium">{stats.tasksByStatus.done}</span>
						</div>
						<div class="w-full bg-gray-200 rounded-full h-2">
							<div
								class="bg-green-500 h-2 rounded-full"
								style="width: {(stats.tasksByStatus.done / stats.totalTasks) * 100}%"
							></div>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
