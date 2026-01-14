import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/svelte';
import MetricCard from './MetricCard.svelte';

describe('MetricCard', () => {
	it('renders with basic props', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Total Sales',
				value: 1500,
				format: 'number'
			}
		});

		expect(getByText('Total Sales')).toBeTruthy();
		expect(getByText('1.5K')).toBeTruthy();
	});

	it('formats currency correctly', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Revenue',
				value: 50000,
				format: 'currency'
			}
		});

		expect(getByText('Revenue')).toBeTruthy();
		expect(getByText('$50,000')).toBeTruthy();
	});

	it('formats percentage correctly', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Conversion Rate',
				value: 15.5,
				format: 'percentage'
			}
		});

		expect(getByText('Conversion Rate')).toBeTruthy();
		expect(getByText('15.5%')).toBeTruthy();
	});

	it('displays positive trend correctly', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Active Users',
				value: 1200,
				trend: 12.5
			}
		});

		expect(getByText('Active Users')).toBeTruthy();
		expect(getByText(/12\.5%/)).toBeTruthy();
	});

	it('displays negative trend correctly', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Bounce Rate',
				value: 45,
				trend: -5.2
			}
		});

		expect(getByText('Bounce Rate')).toBeTruthy();
		expect(getByText(/5\.2%/)).toBeTruthy();
	});

	it('displays neutral trend correctly', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Sessions',
				value: 500,
				trend: 0
			}
		});

		expect(getByText('0%')).toBeTruthy();
	});

	it('renders with icon when provided', () => {
		const mockIcon = {
			name: 'users',
			svg: '<svg><path d="M12 2L2 7l10 5 10-5-10-5z"/></svg>'
		};

		const { container } = render(MetricCard, {
			props: {
				title: 'Users',
				value: 250,
				icon: mockIcon
			}
		});

		expect(container.querySelector('.card')).toBeTruthy();
	});

	it('applies correct color prop', () => {
		const { container } = render(MetricCard, {
			props: {
				title: 'Revenue',
				value: 50000,
				color: 'green'
			}
		});

		expect(container.querySelector('.card')).toBeTruthy();
	});

	it('handles large numbers with correct formatting', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Total Views',
				value: 1234567,
				format: 'number'
			}
		});

		expect(getByText('1.23M')).toBeTruthy();
	});

	it('handles decimal numbers', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Average Score',
				value: 3.14159,
				format: 'number'
			}
		});

		expect(getByText(/3\.\d+/)).toBeTruthy();
	});

	it('formats duration correctly - seconds', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Avg Session',
				value: 45,
				format: 'duration'
			}
		});

		expect(getByText('45s')).toBeTruthy();
	});

	it('formats duration correctly - minutes', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Avg Session',
				value: 150,
				format: 'duration'
			}
		});

		expect(getByText('2.5m')).toBeTruthy();
	});

	it('formats duration correctly - hours', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Avg Session',
				value: 9000,
				format: 'duration'
			}
		});

		expect(getByText('2.5h')).toBeTruthy();
	});

	it('renders without trend when not provided', () => {
		const { container } = render(MetricCard, {
			props: {
				title: 'Count',
				value: 100
			}
		});

		// No trend indicator should be present
		expect(container.querySelector('.card')).toBeTruthy();
	});

	it('handles zero value correctly', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Zero Value',
				value: 0,
				format: 'number'
			}
		});

		expect(getByText('0')).toBeTruthy();
	});

	it('displays trendPeriod text when provided', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Sales',
				value: 1000,
				trend: 10,
				trendPeriod: 'vs last month'
			}
		});

		expect(getByText(/vs last month/)).toBeTruthy();
	});

	it('displays description when provided', () => {
		const { getByText } = render(MetricCard, {
			props: {
				title: 'Sales',
				value: 1000,
				description: 'Total sales this quarter'
			}
		});

		expect(getByText('Total sales this quarter')).toBeTruthy();
	});

	it('handles different color values', () => {
		const colors = ['blue', 'green', 'red', 'purple', 'orange', 'indigo'];

		colors.forEach((color) => {
			const { container } = render(MetricCard, {
				props: {
					title: 'Test',
					value: 100,
					color
				}
			});

			expect(container.querySelector('.card')).toBeTruthy();
		});
	});

	it('shows loading state', () => {
		const { container } = render(MetricCard, {
			props: {
				title: 'Loading',
				value: 100,
				loading: true
			}
		});

		expect(container.querySelector('.animate-pulse')).toBeTruthy();
	});
});
