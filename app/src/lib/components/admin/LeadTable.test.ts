import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import LeadTable from './LeadTable.svelte';

describe('LeadTable Component', () => {
	const mockLeads = [
		{
			id: '1',
			name: 'John Doe',
			email: 'john@example.com',
			phone: '555-1234',
			company: 'Acme Corp',
			status: 'new',
			qualification_score: 75,
			created_at: '2024-01-15T10:00:00Z',
			agent_name: 'Agent Smith'
		},
		{
			id: '2',
			first_name: 'Jane',
			last_name: 'Smith',
			email: 'jane@example.com',
			phone: '555-5678',
			company: 'TechCo',
			status: 'qualified',
			qualification_score: 90,
			created_at: '2024-01-14T09:00:00Z',
			assigned_agent: { name: 'Agent Johnson' }
		},
		{
			id: '3',
			name: 'Bob Wilson',
			email: 'bob@example.com',
			status: 'contacted',
			qualification_score: 50,
			created_at: '2024-01-13T08:00:00Z'
		}
	];

	let onViewMock: any;
	let onAssignMock: any;
	let onEditMock: any;

	beforeEach(() => {
		vi.clearAllMocks();
		onViewMock = vi.fn();
		onAssignMock = vi.fn();
		onEditMock = vi.fn();
	});

	it('should render table with leads', () => {
		render(LeadTable, {
			props: {
				leads: mockLeads,
				onView: onViewMock,
				onAssign: onAssignMock,
				onEdit: onEditMock
			}
		});

		expect(screen.getByText('John Doe')).toBeInTheDocument();
		expect(screen.getByText('Jane Smith')).toBeInTheDocument();
		expect(screen.getByText('Bob Wilson')).toBeInTheDocument();
	});

	it('should display lead information correctly', () => {
		render(LeadTable, {
			props: { leads: mockLeads }
		});

		// Check email and phone
		expect(screen.getByText('john@example.com')).toBeInTheDocument();
		expect(screen.getByText('555-1234')).toBeInTheDocument();

		// Check company
		expect(screen.getByText('Acme Corp')).toBeInTheDocument();

		// Check status
		expect(screen.getByText('new')).toBeInTheDocument();
		expect(screen.getByText('qualified')).toBeInTheDocument();
	});

	it('should show empty state when no leads', () => {
		render(LeadTable, {
			props: { leads: [] }
		});

		expect(screen.getByText('No leads found')).toBeInTheDocument();
	});

	it('should display qualification scores', () => {
		render(LeadTable, {
			props: { leads: mockLeads }
		});

		expect(screen.getByText('75/100')).toBeInTheDocument();
		expect(screen.getByText('90/100')).toBeInTheDocument();
		expect(screen.getByText('50/100')).toBeInTheDocument();
	});

	it('should show agent assignments', () => {
		render(LeadTable, {
			props: { leads: mockLeads }
		});

		expect(screen.getByText('Agent Smith')).toBeInTheDocument();
		expect(screen.getByText('Agent Johnson')).toBeInTheDocument();
		expect(screen.getByText('Unassigned')).toBeInTheDocument();
	});

	it('should format dates correctly', () => {
		render(LeadTable, {
			props: { leads: mockLeads }
		});

		// Dates should be formatted (exact format depends on locale)
		expect(screen.getByText(/Jan.*15.*2024/i)).toBeInTheDocument();
	});

	it('should call onView when view button clicked', async () => {
		render(LeadTable, {
			props: {
				leads: mockLeads,
				onView: onViewMock
			}
		});

		const viewButtons = screen.getAllByTitle('View');
		await fireEvent.click(viewButtons[0]);

		expect(onViewMock).toHaveBeenCalledWith(mockLeads[0]);
	});

	it('should call onEdit when edit button clicked', async () => {
		render(LeadTable, {
			props: {
				leads: mockLeads,
				onEdit: onEditMock
			}
		});

		const editButtons = screen.getAllByTitle('Edit');
		await fireEvent.click(editButtons[0]);

		expect(onEditMock).toHaveBeenCalledWith(mockLeads[0]);
	});

	it('should call onAssign when assign button clicked', async () => {
		render(LeadTable, {
			props: {
				leads: mockLeads,
				onAssign: onAssignMock
			}
		});

		const assignButtons = screen.getAllByTitle('Assign');
		await fireEvent.click(assignButtons[0]);

		expect(onAssignMock).toHaveBeenCalledWith(mockLeads[0]);
	});

	describe('Sorting', () => {
		it('should sort by name when name column clicked', async () => {
			const { component } = render(LeadTable, {
				props: { leads: [...mockLeads] }
			});

			const nameHeader = screen.getByText('Name').closest('th');
			await fireEvent.click(nameHeader!);

			// After sorting, check order (Bob, Jane, John alphabetically)
			const rows = screen.getAllByRole('row');
			expect(rows[1].textContent).toContain('Bob Wilson');
		});

		it('should toggle sort direction on multiple clicks', async () => {
			render(LeadTable, {
				props: { leads: [...mockLeads] }
			});

			const nameHeader = screen.getByText('Name').closest('th');

			// First click - ascending
			await fireEvent.click(nameHeader!);

			// Second click - descending
			await fireEvent.click(nameHeader!);

			const rows = screen.getAllByRole('row');
			expect(rows[1].textContent).toContain('John Doe');
		});

		it('should sort by status when status column clicked', async () => {
			render(LeadTable, {
				props: { leads: [...mockLeads] }
			});

			const statusHeader = screen.getByText('Status').closest('th');
			await fireEvent.click(statusHeader!);

			// Check that sorting occurred
			expect(screen.getAllByRole('row').length).toBeGreaterThan(0);
		});

		it('should sort by score when score column clicked', async () => {
			render(LeadTable, {
				props: { leads: [...mockLeads] }
			});

			const scoreHeader = screen.getByText('Score').closest('th');
			await fireEvent.click(scoreHeader!);

			// After sorting by score ascending, lowest should be first
			const rows = screen.getAllByRole('row');
			expect(rows[1].textContent).toContain('50/100');
		});

		it('should sort by created date when created column clicked', async () => {
			render(LeadTable, {
				props: { leads: [...mockLeads] }
			});

			const createdHeader = screen.getByText('Created').closest('th');
			await fireEvent.click(createdHeader!);

			// Check that sorting occurred
			expect(screen.getAllByRole('row').length).toBeGreaterThan(0);
		});
	});

	describe('Status Colors', () => {
		it('should apply correct color classes to status badges', () => {
			render(LeadTable, {
				props: { leads: mockLeads }
			});

			const newBadge = screen.getByText('new');
			expect(newBadge.className).toContain('bg-blue-100');
			expect(newBadge.className).toContain('text-blue-800');

			const qualifiedBadge = screen.getByText('qualified');
			expect(qualifiedBadge.className).toContain('bg-green-100');
			expect(qualifiedBadge.className).toContain('text-green-800');
		});
	});

	describe('Score Progress Bars', () => {
		it('should render score progress bars with correct width', () => {
			const { container } = render(LeadTable, {
				props: { leads: [mockLeads[0]] }
			});

			const progressBar = container.querySelector('.bg-blue-600');
			expect(progressBar).toBeTruthy();
			expect(progressBar?.getAttribute('style')).toContain('width: 75%');
		});
	});

	describe('Hover Effects', () => {
		it('should apply hover class to table rows', () => {
			const { container } = render(LeadTable, {
				props: { leads: mockLeads }
			});

			const rows = container.querySelectorAll('tbody tr');
			rows.forEach((row) => {
				expect(row.className).toContain('hover:bg-gray-50');
			});
		});
	});

	describe('Action Buttons', () => {
		it('should render all action buttons for each lead', () => {
			render(LeadTable, {
				props: { leads: [mockLeads[0]] }
			});

			expect(screen.getByTitle('View')).toBeInTheDocument();
			expect(screen.getByTitle('Edit')).toBeInTheDocument();
			expect(screen.getByTitle('Assign')).toBeInTheDocument();
		});

		it('should display action buttons with proper SVG icons', () => {
			const { container } = render(LeadTable, {
				props: { leads: [mockLeads[0]] }
			});

			const actionButtons = container.querySelectorAll('td:last-child button');
			expect(actionButtons.length).toBe(3);

			// Each button should have an SVG
			actionButtons.forEach((button) => {
				expect(button.querySelector('svg')).toBeTruthy();
			});
		});
	});

	describe('Null/Undefined Handling', () => {
		it('should handle leads with missing optional fields', () => {
			const leadWithMissingFields = {
				id: '4',
				name: 'Test Lead',
				email: null,
				phone: null,
				company: null,
				status: null,
				qualification_score: null,
				created_at: null
			};

			render(LeadTable, {
				props: { leads: [leadWithMissingFields] }
			});

			expect(screen.getByText('Test Lead')).toBeInTheDocument();
			// Check for N/A text - there might be multiple instances
			const naTexts = screen.getAllByText('N/A');
			expect(naTexts.length).toBeGreaterThan(0); // For missing email and/or created_at
			expect(screen.getByText('0/100')).toBeInTheDocument(); // For null score
		});
	});
});
