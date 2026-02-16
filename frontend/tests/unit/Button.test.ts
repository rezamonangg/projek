import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/svelte';
import Button from '$lib/components/common/Button.svelte';

describe('Button', () => {
	it('renders with default props', () => {
		const { getByRole } = render(Button, { children: 'Click me' });
		expect(getByRole('button')).toBeInTheDocument();
	});

	it('renders with primary variant', () => {
		const { getByRole } = render(Button, { variant: 'primary', children: 'Click me' });
		const button = getByRole('button');
		expect(button).toHaveClass('bg-blue-600');
	});

	it('renders with secondary variant', () => {
		const { getByRole } = render(Button, { variant: 'secondary', children: 'Click me' });
		const button = getByRole('button');
		expect(button).toHaveClass('bg-gray-200');
	});

	it('shows loading state', () => {
		const { getByRole } = render(Button, { loading: true, children: 'Click me' });
		const button = getByRole('button');
		expect(button).toHaveAttribute('aria-busy', 'true');
	});

	it('is disabled when disabled prop is true', () => {
		const { getByRole } = render(Button, { disabled: true, children: 'Click me' });
		const button = getByRole('button');
		expect(button).toBeDisabled();
	});
});
