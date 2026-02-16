<script lang="ts">
	type InputType = 'text' | 'email' | 'password' | 'number' | 'tel' | 'url' | 'search';

	interface Props {
		type?: InputType;
		name?: string;
		value?: string;
		placeholder?: string;
		disabled?: boolean;
		readonly?: boolean;
		required?: boolean;
		error?: string;
		label?: string;
		id?: string;
		class?: string;
		oninput?: (e: Event) => void;
		onchange?: (e: Event) => void;
		onblur?: (e: FocusEvent) => void;
		onfocus?: (e: FocusEvent) => void;
	}

	let {
		type = 'text',
		name,
		value = $bindable(''),
		placeholder,
		disabled = false,
		readonly = false,
		required = false,
		error,
		label,
		id,
		class: className = '',
		oninput,
		onchange,
		onblur,
		onfocus
	}: Props = $props();

	let inputId = $derived(id || name || crypto.randomUUID());
</script>

<div class="w-full">
	{#if label}
		<label for={inputId} class="block text-sm font-medium text-gray-700 mb-1">
			{label}
			{#if required}
				<span class="text-red-500">*</span>
			{/if}
		</label>
	{/if}
	<input
		{id}
		{name}
		{type}
		bind:value
		{placeholder}
		{disabled}
		{readonly}
		{required}
		{oninput}
		{onchange}
		{onblur}
		{onfocus}
		class="w-full px-3 py-2 border rounded-lg shadow-sm transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed {error
			? 'border-red-500 focus:ring-red-500 focus:border-red-500'
			: 'border-gray-300'} {className}"
		aria-invalid={error ? 'true' : undefined}
		aria-describedby={error ? `${inputId}-error` : undefined}
	/>
	{#if error}
		<p id="{inputId}-error" class="mt-1 text-sm text-red-600">
			{error}
		</p>
	{/if}
</div>
