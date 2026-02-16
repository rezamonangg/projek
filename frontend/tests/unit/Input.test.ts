import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/svelte';
import Input from '$lib/components/common/Input.svelte';

describe('Input', () => {
	it('renders with label', () => {
		const { getByLabelText } = render(Input, { label: 'Email', name: 'email' });
		expect(getByLabelText('Email')).toBeInTheDocument();
	});

	it('renders with placeholder', () => {
		const { getByPlaceholderText } = render(Input, {
			placeholder: 'Enter email',
			name: 'email'
		});
		expect(getByPlaceholderText('Enter email')).toBeInTheDocument();
	});

	it('shows error message', () => {
		const { getByText } = render(Input, {
			label: 'Email',
			name: 'email',
			error: 'Invalid email'
		});
		expect(getByText('Invalid email')).toBeInTheDocument();
	});

	it('shows required indicator', () => {
		const { getByText } = render(Input, {
			label: 'Email',
			name: 'email',
			required: true
		});
		expect(getByText('*')).toBeInTheDocument();
	});

	it('is disabled when disabled prop is true', () => {
		const { getByRole } = render(Input, {
			name: 'email',
			disabled: true
		});
		expect(getByRole('textbox')).toBeDisabled();
	});
});
