<script module lang="ts">
	// Declare Facebook Pixel for TypeScript
	declare const fbq: ((command: string, event: string) => void) | undefined;
</script>

<script lang="ts">
	import {
		Calendar,
		Clock,
		ChevronLeft,
		ChevronRight,
		Video,
		CheckCircle,
		X,
		Loader
	} from 'lucide-svelte';
	import { trackEvent, trackConversion } from '$lib/analytics';
	import { trackFunnelStep } from '$lib/funnel';
	import { onMount } from 'svelte';
	import * as api from '$lib/api';

	interface Props {
		leadId?: string;
		leadEmail?: string;
		leadName?: string;
		mode?: 'embed' | 'custom';
		calendarEmbedUrl?: string;
	}

	let {
		leadId = '',
		leadEmail = '',
		leadName = '',
		mode = 'custom',
		calendarEmbedUrl = ''
	}: Props = $props();

	// State management
	let currentDate = $state(new Date());
	let selectedDate = $state<Date | null>(null);
	let selectedTime = $state('');
	let timezone = $state(Intl.DateTimeFormat().resolvedOptions().timeZone);
	let phone = $state('');
	let notes = $state('');
	let isSubmitting = $state(false);
	let bookingComplete = $state(false);
	let bookingDetails = $state<any>(null);
	let error = $state('');
	let availabilityCache = $state<Map<string, any>>(new Map());
	let dateAvailabilityCache = $state<Map<string, boolean>>(new Map());
	let isLoadingAvailability = $state(false);
	let isLoadingCalendar = $state(true);
	let currentView = $state<'calendar' | 'timeslots' | 'form'>('calendar');

	// On mount, navigate to the first available date and pre-fetch availability
	onMount(() => {
		navigateToFirstAvailableDate();
		prefetchMonthAvailability();
	});

	// Find and navigate to the first available date
	function navigateToFirstAvailableDate() {
		const today = new Date();
		today.setHours(0, 0, 0, 0);

		// Check next 21 days (3 weeks) for availability
		for (let i = 0; i < 21; i++) {
			const checkDate = new Date(today);
			checkDate.setDate(checkDate.getDate() + i);

			// Check if this date is available (M-Th only, not weekend/Friday)
			const day = checkDate.getDay();
			if (day !== 0 && day !== 5 && day !== 6) {
				// Found an available date - make sure currentDate shows this month
				if (
					checkDate.getMonth() !== currentDate.getMonth() ||
					checkDate.getFullYear() !== currentDate.getFullYear()
				) {
					currentDate = new Date(checkDate.getFullYear(), checkDate.getMonth(), 1);
				}
				return;
			}
		}
	}

	// Pre-fetch availability for all potentially available dates in the visible range
	async function prefetchMonthAvailability() {
		isLoadingCalendar = true;
		const today = new Date();
		today.setHours(0, 0, 0, 0);

		const maxDate = new Date(today);
		maxDate.setDate(maxDate.getDate() + 21);

		// Get all M-Th dates within the next 21 days
		const datesToCheck: Date[] = [];
		for (let i = 0; i < 21; i++) {
			const checkDate = new Date(today);
			checkDate.setDate(checkDate.getDate() + i);
			const day = checkDate.getDay();
			// Only check M-Th (1-4)
			if (day >= 1 && day <= 4) {
				datesToCheck.push(checkDate);
			}
		}

		// Fetch availability for all dates in parallel
		const promises = datesToCheck.map(async (date) => {
			const dateStr = formatDateForAPI(date);
			try {
				const data = await api.getAvailability({ date: dateStr, timezone });
				availabilityCache.set(dateStr, data);
				// Check if there are any available slots
				const hasAvailableSlots =
					data.availability?.some((slot: { available: boolean }) => slot.available) ?? false;
				dateAvailabilityCache.set(dateStr, hasAvailableSlots);
			} catch (err) {
				console.error('Error checking availability for', dateStr, err);
				dateAvailabilityCache.set(dateStr, false);
			}
		});

		await Promise.all(promises);
		// Force reactivity update
		dateAvailabilityCache = new Map(dateAvailabilityCache);
		isLoadingCalendar = false;
	}

	// Common timezones
	const timezones = [
		'America/New_York',
		'America/Chicago',
		'America/Denver',
		'America/Los_Angeles',
		'America/Phoenix',
		'America/Anchorage',
		'Pacific/Honolulu',
		'Europe/London',
		'Europe/Paris',
		'Europe/Berlin',
		'Asia/Tokyo',
		'Asia/Shanghai',
		'Asia/Dubai',
		'Australia/Sydney',
		'Pacific/Auckland'
	];

	// Get days in month
	function getDaysInMonth(date: Date): Date[] {
		const year = date.getFullYear();
		const month = date.getMonth();
		const firstDay = new Date(year, month, 1);
		const lastDay = new Date(year, month + 1, 0);
		const daysInMonth = lastDay.getDate();

		const days: Date[] = [];

		// Add padding days from previous month
		const firstDayOfWeek = firstDay.getDay();
		for (let i = firstDayOfWeek - 1; i >= 0; i--) {
			const prevDate = new Date(year, month, -i);
			days.push(prevDate);
		}

		// Add days of current month
		for (let day = 1; day <= daysInMonth; day++) {
			days.push(new Date(year, month, day));
		}

		// Add padding days from next month
		const remainingDays = 7 - (days.length % 7);
		if (remainingDays < 7) {
			for (let i = 1; i <= remainingDays; i++) {
				days.push(new Date(year, month + 1, i));
			}
		}

		return days;
	}

	// Navigation
	function previousMonth() {
		currentDate = new Date(currentDate.getFullYear(), currentDate.getMonth() - 1, 1);
	}

	function nextMonth() {
		currentDate = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 1);
	}

	// Check if date is available (basic checks - day of week, past, future)
	function isDatePotentiallyAvailable(date: Date): boolean {
		const today = new Date();
		today.setHours(0, 0, 0, 0);

		// Check if past
		if (date < today) return false;

		// Check if more than 21 days (3 weeks) in the future
		const maxDate = new Date(today);
		maxDate.setDate(maxDate.getDate() + 21);
		if (date > maxDate) return false;

		// Check if weekend or Friday (only M-Th available)
		// Sunday = 0, Monday = 1, Tuesday = 2, Wednesday = 3, Thursday = 4, Friday = 5, Saturday = 6
		const day = date.getDay();
		if (day === 0 || day === 5 || day === 6) return false;

		// Check if same month
		if (date.getMonth() !== currentDate.getMonth()) return false;

		return true;
	}

	// Check if date has actual available time slots
	function isDateAvailable(date: Date): boolean {
		// First check basic availability
		if (!isDatePotentiallyAvailable(date)) return false;

		// If still loading, show as potentially available
		if (isLoadingCalendar) return true;

		// Check the pre-fetched availability cache
		const dateStr = formatDateForAPI(date);
		const hasSlots = dateAvailabilityCache.get(dateStr);

		// If we have cached data, use it; otherwise assume unavailable
		return hasSlots ?? false;
	}

	// Check if date has available slots
	async function checkDateAvailability(date: Date): Promise<number> {
		const dateStr = formatDateForAPI(date);

		// Check cache first
		if (availabilityCache.has(dateStr)) {
			const cached = availabilityCache.get(dateStr);
			return cached.availability.filter((slot: any) => slot.available).length;
		}

		try {
			const data = await api.getAvailability({ date: dateStr, timezone });
			availabilityCache.set(dateStr, data);
			return data.availability.filter((slot: { available: boolean }) => slot.available).length;
		} catch (err) {
			console.error('Error checking availability:', err);
		}

		return 0;
	}

	// Format date for API
	function formatDateForAPI(date: Date): string {
		const year = date.getFullYear();
		const month = String(date.getMonth() + 1).padStart(2, '0');
		const day = String(date.getDate()).padStart(2, '0');
		return `${year}-${month}-${day}`;
	}

	// Format date for display
	function formatDateDisplay(date: Date): string {
		return date.toLocaleDateString('en-US', {
			weekday: 'long',
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}

	// Format time from 24-hour (HH:MM) to 12-hour AM/PM format in user's timezone
	function formatTimeDisplay(time24: string): string {
		// Parse the 24-hour time (e.g., "09:00", "13:00")
		const [hours, minutes] = time24.split(':').map(Number);

		// Create a date object for today with the given time
		// The time from the API is in the selected timezone
		const period = hours >= 12 ? 'PM' : 'AM';
		const hours12 = hours % 12 || 12;

		return `${hours12}:${minutes.toString().padStart(2, '0')} ${period}`;
	}

	// Handle date selection
	async function selectDate(date: Date) {
		if (!isDateAvailable(date)) return;

		selectedDate = date;
		currentView = 'timeslots';
		isLoadingAvailability = true;
		error = '';

		// Track calendar date selection
		trackEvent('calendar_date_selected', {
			date: formatDateForAPI(date),
			timezone
		});

		// Track funnel step: calendar selection
		trackFunnelStep('calendar');

		const dateStr = formatDateForAPI(date);

		try {
			const data = await api.getAvailability({ date: dateStr, timezone });
			availabilityCache.set(dateStr, data);
		} catch (err) {
			error = 'Failed to load available times. Please try again.';
		} finally {
			isLoadingAvailability = false;
		}
	}

	// Handle time selection
	function selectTime(time: string) {
		selectedTime = time;
		currentView = 'form';

		// Track time slot selection
		trackEvent('calendar_time_selected', {
			date: selectedDate ? formatDateForAPI(selectedDate) : '',
			time,
			timezone
		});
	}

	// Handle booking submission
	async function handleBooking() {
		if (!selectedDate || !selectedTime) {
			error = 'Please select a date and time';
			return;
		}

		isSubmitting = true;
		error = '';

		try {
			const data = await api.bookMeeting({
				lead_id: leadId,
				attendee_name: leadName,
				attendee_email: leadEmail,
				attendee_phone: phone || undefined,
				meeting_date: formatDateForAPI(selectedDate),
				meeting_time: selectedTime,
				timezone: timezone,
				notes: notes || undefined
			});

			bookingDetails = data;
			bookingComplete = true;

			// Facebook Pixel conversion event
			if (typeof fbq !== 'undefined') {
				fbq('track', 'Schedule');
			}

			// Track successful booking conversion
			trackConversion('meeting_booked', {
				date: selectedDate ? formatDateForAPI(selectedDate) : '',
				time: selectedTime,
				timezone,
				meeting_id: data.meeting_id
			});

			// Track funnel step: completion
			trackFunnelStep('completion', undefined, leadId);
		} catch (err: any) {
			error = err.message || 'Failed to book meeting. Please try again.';
		} finally {
			isSubmitting = false;
		}
	}

	// Go back to calendar
	function backToCalendar() {
		selectedDate = null;
		selectedTime = '';
		currentView = 'calendar';
		error = '';
	}

	// Go back to time slots
	function backToTimeSlots() {
		selectedTime = '';
		currentView = 'timeslots';
		error = '';
	}

	// Generate .ics calendar file
	function generateICSFile(): string {
		if (!bookingDetails || !selectedDate || !selectedTime) return '';

		const startDate = new Date(`${formatDateForAPI(selectedDate)}T${selectedTime}:00`);
		const endDate = new Date(startDate.getTime() + 30 * 60000); // 30 minutes

		const formatICSDate = (date: Date) => {
			return date.toISOString().replace(/[-:]/g, '').split('.')[0] + 'Z';
		};

		const icsContent = [
			'BEGIN:VCALENDAR',
			'VERSION:2.0',
			'PRODID:-//Outlet//Calendar Booking//EN',
			'BEGIN:VEVENT',
			`UID:${bookingDetails.meeting_id}@outlet.sh`,
			`DTSTAMP:${formatICSDate(new Date())}`,
			`DTSTART:${formatICSDate(startDate)}`,
			`DTEND:${formatICSDate(endDate)}`,
			`SUMMARY:Outlet Consultation`,
			`DESCRIPTION:Initial consultation with Outlet`,
			bookingDetails.zoom_join_url ? `URL:${bookingDetails.zoom_join_url}` : '',
			'STATUS:CONFIRMED',
			'END:VEVENT',
			'END:VCALENDAR'
		]
			.filter(Boolean)
			.join('\r\n');

		return icsContent;
	}

	function downloadICS() {
		const icsContent = generateICSFile();
		const blob = new Blob([icsContent], { type: 'text/calendar;charset=utf-8' });
		const link = document.createElement('a');
		link.href = URL.createObjectURL(blob);
		link.download = 'outlet-consultation.ics';
		link.click();
	}

	// Get current month/year display
	const monthYearDisplay = $derived(
		currentDate.toLocaleDateString('en-US', {
			month: 'long',
			year: 'numeric'
		})
	);

	// Get calendar days
	const calendarDays = $derived(getDaysInMonth(currentDate));

	// Get available time slots for selected date (filtering past slots for today)
	const availableTimeSlots = $derived.by(() => {
		if (!selectedDate) return [];

		const dateStr = formatDateForAPI(selectedDate);
		const cached = availabilityCache.get(dateStr);

		if (!cached) return [];

		const slots = cached.availability || [];

		// If selected date is today, filter out time slots that have already passed
		const today = new Date();
		const isToday = selectedDate.toDateString() === today.toDateString();

		if (!isToday) return slots;

		const currentHour = today.getHours();
		const currentMinute = today.getMinutes();

		return slots.map((slot: any) => {
			const [slotHour, slotMinute] = slot.time.split(':').map(Number);

			// Mark as unavailable if the slot time has passed
			// Add 30 min buffer - don't allow booking within 30 minutes
			const slotTotalMinutes = slotHour * 60 + slotMinute;
			const currentTotalMinutes = currentHour * 60 + currentMinute + 30; // 30 min buffer

			if (slotTotalMinutes <= currentTotalMinutes) {
				return { ...slot, available: false };
			}

			return slot;
		});
	});
</script>

<div class="calendar-booking">
	{#if mode === 'embed' && calendarEmbedUrl}
		<!-- Embedded Calendar -->
		<div class="embed-container">
			<iframe
				src={calendarEmbedUrl}
				width="100%"
				height="700"
				frameborder="0"
				title="Schedule a meeting"
			></iframe>
		</div>
	{:else if bookingComplete}
		<!-- Booking Confirmation -->
		<div class="confirmation-card">
			<div class="confirmation-icon">
				<CheckCircle class="w-16 h-16 text-green-600" />
			</div>

			<h2 class="confirmation-title">Meeting Scheduled!</h2>
			<p class="confirmation-subtitle">Your consultation has been confirmed</p>

			<div class="confirmation-details">
				<div class="detail-item">
					<Calendar class="w-5 h-5 text-blue-600" />
					<div>
						<div class="detail-label">Date & Time</div>
						<div class="detail-value">
							{selectedDate ? formatDateDisplay(selectedDate) : ''} at {selectedTime
								? formatTimeDisplay(selectedTime)
								: ''}
						</div>
						<div class="detail-timezone">{timezone}</div>
					</div>
				</div>

				<div class="detail-item">
					<Clock class="w-5 h-5 text-blue-600" />
					<div>
						<div class="detail-label">Duration</div>
						<div class="detail-value">30 minutes</div>
					</div>
				</div>

				{#if bookingDetails?.zoom_join_url}
					<div class="detail-item">
						<Video class="w-5 h-5 text-blue-600" />
						<div>
							<div class="detail-label">Meeting Link</div>
							<a href={bookingDetails.zoom_join_url} class="detail-link" target="_blank"
								>Join Google Meet</a
							>
						</div>
					</div>
				{/if}
			</div>

			<div class="confirmation-message">
				<p>
					A calendar invite has been sent to <strong>{leadEmail}</strong>. You'll receive a reminder
					email 24 hours before the meeting.
				</p>
			</div>

			<div class="button-group">
				<button onclick={downloadICS} class="secondary-button">
					<Calendar class="w-4 h-4" />
					Download .ics File
				</button>
				<a href="/" class="primary-button">Return to Homepage</a>
			</div>
		</div>
	{:else}
		<!-- Custom Calendar Interface -->
		<div class="booking-card">
			{#if currentView === 'calendar'}
				<!-- Month Calendar View -->
				<div class="calendar-header">
					<h2 class="calendar-title">Select a Date</h2>
					<p class="calendar-subtitle">Choose a day for your consultation</p>

					<!-- Timezone Selector -->
					<div class="timezone-selector">
						<Clock class="w-4 h-4" />
						<select bind:value={timezone} class="timezone-select">
							{#each timezones as tz}
								<option value={tz}>{tz.replace(/_/g, ' ')}</option>
							{/each}
						</select>
					</div>
				</div>

				<!-- Month Navigation -->
				<div class="month-navigation">
					<button onclick={previousMonth} class="nav-button">
						<ChevronLeft class="w-5 h-5" />
					</button>
					<h3 class="month-title">{monthYearDisplay}</h3>
					<button onclick={nextMonth} class="nav-button">
						<ChevronRight class="w-5 h-5" />
					</button>
				</div>

				<!-- Calendar Grid -->
				<div class="calendar-grid-container">
					{#if isLoadingCalendar}
						<div class="calendar-loading-overlay">
							<Loader class="w-6 h-6 animate-spin text-blue-600" />
							<span class="text-sm text-gray-600">Checking availability...</span>
						</div>
					{/if}
					<div class="calendar-grid" class:loading={isLoadingCalendar}>
						<!-- Day headers -->
						{#each ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'] as day}
							<div class="day-header">{day}</div>
						{/each}

						<!-- Calendar days -->
						{#each calendarDays as day}
							{@const isAvailable = isDateAvailable(day)}
							{@const isToday = day.toDateString() === new Date().toDateString()}
							{@const isCurrentMonth = day.getMonth() === currentDate.getMonth()}
							<button
								onclick={() => selectDate(day)}
								disabled={!isAvailable || isLoadingCalendar}
								class="calendar-day"
								class:available={isAvailable}
								class:unavailable={!isAvailable}
								class:today={isToday}
								class:other-month={!isCurrentMonth}
							>
								<span class="day-number">{day.getDate()}</span>
								{#if isAvailable && !isLoadingCalendar}
									<span class="availability-indicator"></span>
								{/if}
							</button>
						{/each}
					</div>
				</div>

				<div class="calendar-legend">
					<div class="legend-item">
						<span class="legend-dot available"></span>
						<span>Available</span>
					</div>
					<div class="legend-item">
						<span class="legend-dot unavailable"></span>
						<span>Unavailable</span>
					</div>
				</div>
			{:else if currentView === 'timeslots'}
				<!-- Time Slot Selection -->
				<div class="timeslots-header">
					<button onclick={backToCalendar} class="back-button">
						<ChevronLeft class="w-4 h-4" />
						Back to Calendar
					</button>
					<h2 class="timeslots-title">Select a Time</h2>
					<p class="timeslots-subtitle">
						{selectedDate ? formatDateDisplay(selectedDate) : ''}
					</p>
					<p class="timezone-info">All times shown in {timezone}</p>
				</div>

				{#if error}
					<div class="error-message">
						<X class="w-5 h-5" />
						<span>{error}</span>
					</div>
				{/if}

				{#if isLoadingAvailability}
					<div class="loading-container">
						<Loader class="w-8 h-8 animate-spin text-blue-600" />
						<p>Loading available times...</p>
					</div>
				{:else if availableTimeSlots.length === 0}
					<div class="no-slots-message">
						<p>No available time slots for this date.</p>
						<button onclick={backToCalendar} class="secondary-button"> Choose Another Date </button>
					</div>
				{:else}
					<div class="timeslots-grid">
						{#each availableTimeSlots as slot}
							<button
								onclick={() => selectTime(slot.time)}
								disabled={!slot.available}
								class="time-slot"
								class:available={slot.available}
								class:unavailable={!slot.available}
							>
								{formatTimeDisplay(slot.time)}
							</button>
						{/each}
					</div>
				{/if}
			{:else if currentView === 'form'}
				<!-- Booking Form -->
				<div class="form-header">
					<button onclick={backToTimeSlots} class="back-button">
						<ChevronLeft class="w-4 h-4" />
						Back to Times
					</button>
					<h2 class="form-title">Confirm Your Booking</h2>
					<p class="form-subtitle">Review your details and confirm</p>
				</div>

				{#if error}
					<div class="error-message">
						<X class="w-5 h-5" />
						<span>{error}</span>
					</div>
				{/if}

				<!-- Booking Summary -->
				<div class="booking-summary">
					<div class="summary-item">
						<Calendar class="w-5 h-5 text-gray-500" />
						<div>
							<div class="summary-label">Date & Time</div>
							<div class="summary-value">
								{selectedDate ? formatDateDisplay(selectedDate) : ''} at {selectedTime
									? formatTimeDisplay(selectedTime)
									: ''}
							</div>
							<div class="summary-timezone">{timezone}</div>
						</div>
					</div>
					<div class="summary-item">
						<Clock class="w-5 h-5 text-gray-500" />
						<div>
							<div class="summary-label">Duration</div>
							<div class="summary-value">30 minutes</div>
						</div>
					</div>
				</div>

				<!-- Contact Form -->
				<form
					class="booking-form"
					onsubmit={(e) => {
						e.preventDefault();
						handleBooking();
					}}
				>
					<div class="form-group">
						<label for="name" class="form-label">Name</label>
						<input type="text" id="name" value={leadName} readonly class="form-input readonly" />
					</div>

					<div class="form-group">
						<label for="email" class="form-label">Email</label>
						<input type="email" id="email" value={leadEmail} readonly class="form-input readonly" />
					</div>

					<div class="form-group">
						<label for="phone" class="form-label">Phone Number (Optional)</label>
						<input
							type="tel"
							id="phone"
							bind:value={phone}
							placeholder="+1 (555) 123-4567"
							class="form-input"
						/>
					</div>

					<div class="form-group">
						<label for="notes" class="form-label">Additional Notes (Optional)</label>
						<textarea
							id="notes"
							bind:value={notes}
							rows="3"
							placeholder="Any specific topics you'd like to discuss?"
							class="form-input"
						></textarea>
					</div>

					<button type="submit" disabled={isSubmitting} class="submit-button">
						{#if isSubmitting}
							<Loader class="w-5 h-5 animate-spin" />
							Scheduling...
						{:else}
							Confirm Booking
						{/if}
					</button>
				</form>
			{/if}
		</div>
	{/if}
</div>

<style>
	@reference "../../app.css";

	.calendar-booking {
		@apply w-full;
	}

	.embed-container {
		@apply w-full overflow-hidden rounded-2xl border border-gray-200 shadow-lg;
	}

	.booking-card {
		@apply rounded-2xl border border-gray-200 bg-white p-6 shadow-lg md:p-8;
	}

	/* Calendar View Styles */
	.calendar-header {
		@apply mb-6;
	}

	.calendar-title {
		@apply mb-1 text-2xl font-bold text-gray-900;
	}

	.calendar-subtitle {
		@apply mb-4 text-gray-600;
	}

	.timezone-selector {
		@apply mt-4 flex items-center gap-2 text-sm text-gray-600;
	}

	.timezone-select {
		@apply rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-200;
	}

	.month-navigation {
		@apply mb-6 flex items-center justify-between border-b border-gray-200 pb-4;
	}

	.nav-button {
		@apply rounded-lg p-2 transition-colors hover:bg-gray-100;
	}

	.month-title {
		@apply text-xl font-semibold text-gray-900;
	}

	.calendar-grid-container {
		@apply relative;
	}

	.calendar-loading-overlay {
		@apply absolute inset-0 z-10 flex flex-col items-center justify-center gap-2 bg-white/80;
	}

	.calendar-grid {
		@apply grid grid-cols-7 gap-2;
	}

	.calendar-grid.loading {
		@apply opacity-50;
	}

	.day-header {
		@apply py-2 text-center text-sm font-semibold text-gray-600;
	}

	.calendar-day {
		@apply relative flex aspect-square flex-col items-center justify-center rounded-lg border-2 transition-all;
	}

	.calendar-day.available {
		@apply cursor-pointer border-green-200 bg-green-50 hover:border-green-300 hover:bg-green-100;
	}

	.calendar-day.unavailable {
		@apply cursor-not-allowed border-gray-200 bg-gray-50 text-gray-400;
	}

	.calendar-day.today {
		@apply ring-2 ring-blue-500 ring-offset-2;
	}

	.calendar-day.other-month {
		@apply opacity-40;
	}

	.day-number {
		@apply text-base font-medium;
	}

	.availability-indicator {
		@apply absolute bottom-1 h-1.5 w-1.5 rounded-full bg-green-500;
	}

	.calendar-legend {
		@apply mt-6 flex items-center gap-4 border-t border-gray-200 pt-4 text-sm;
	}

	.legend-item {
		@apply flex items-center gap-2;
	}

	.legend-dot {
		@apply h-3 w-3 rounded-full;
	}

	.legend-dot.available {
		@apply bg-green-500;
	}

	.legend-dot.unavailable {
		@apply bg-gray-300;
	}

	/* Time Slots View Styles */
	.timeslots-header {
		@apply mb-6;
	}

	.back-button {
		@apply mb-4 flex items-center gap-1 text-blue-600 transition-colors hover:text-blue-700;
	}

	.timeslots-title {
		@apply mb-1 text-2xl font-bold text-gray-900;
	}

	.timeslots-subtitle {
		@apply mb-1 text-gray-600;
	}

	.timezone-info {
		@apply text-sm text-gray-500;
	}

	.timeslots-grid {
		@apply grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4;
	}

	.time-slot {
		@apply rounded-lg border-2 px-4 py-3 font-medium transition-all;
	}

	.time-slot.available {
		@apply border-blue-200 bg-blue-50 text-blue-700 hover:border-blue-300 hover:bg-blue-100;
	}

	.time-slot.unavailable {
		@apply cursor-not-allowed border-gray-200 bg-gray-50 text-gray-400;
	}

	.loading-container {
		@apply flex flex-col items-center justify-center py-12 text-gray-600;
	}

	.no-slots-message {
		@apply py-12 text-center;
	}

	/* Form View Styles */
	.form-header {
		@apply mb-6;
	}

	.form-title {
		@apply mb-1 text-2xl font-bold text-gray-900;
	}

	.form-subtitle {
		@apply text-gray-600;
	}

	.booking-summary {
		@apply mb-6 space-y-3 rounded-lg border border-blue-200 bg-blue-50 p-4;
	}

	.summary-item {
		@apply flex items-start gap-3;
	}

	.summary-label {
		@apply text-sm text-gray-600;
	}

	.summary-value {
		@apply font-medium text-gray-900;
	}

	.summary-timezone {
		@apply text-sm text-gray-500;
	}

	.booking-form {
		@apply space-y-4;
	}

	.form-group {
		@apply space-y-2;
	}

	.form-label {
		@apply block text-sm font-semibold text-gray-700;
	}

	.form-input {
		@apply w-full rounded-lg border border-gray-300 px-4 py-3 transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-200;
	}

	.form-input.readonly {
		@apply cursor-not-allowed bg-gray-50 text-gray-600;
	}

	.submit-button {
		@apply flex w-full items-center justify-center gap-2 rounded-lg bg-blue-600 px-6 py-4 font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-gray-400;
	}

	/* Confirmation Styles */
	.confirmation-card {
		@apply rounded-2xl border border-gray-200 bg-white p-8 text-center shadow-lg;
	}

	.confirmation-icon {
		@apply mb-6 flex justify-center;
	}

	.confirmation-title {
		@apply mb-2 text-3xl font-bold text-gray-900;
	}

	.confirmation-subtitle {
		@apply mb-8 text-gray-600;
	}

	.confirmation-details {
		@apply mb-8 space-y-4 text-left;
	}

	.detail-item {
		@apply flex items-start gap-4 rounded-lg bg-gray-50 p-4;
	}

	.detail-label {
		@apply mb-1 text-sm text-gray-600;
	}

	.detail-value {
		@apply font-medium text-gray-900;
	}

	.detail-timezone {
		@apply text-sm text-gray-500;
	}

	.detail-link {
		@apply font-medium text-blue-600 underline hover:text-blue-700;
	}

	.confirmation-message {
		@apply mb-6 rounded-lg border border-blue-200 bg-blue-50 p-4;
	}

	.confirmation-message p {
		@apply text-gray-700;
	}

	.button-group {
		@apply flex flex-col justify-center gap-3 sm:flex-row;
	}

	.primary-button {
		@apply inline-flex items-center justify-center rounded-lg bg-blue-600 px-6 py-3 font-semibold text-white transition-colors hover:bg-blue-700;
	}

	.secondary-button {
		@apply inline-flex items-center justify-center gap-2 rounded-lg border border-gray-300 bg-white px-6 py-3 font-semibold text-gray-700 transition-colors hover:bg-gray-50;
	}

	.error-message {
		@apply mb-4 flex items-center gap-2 rounded-lg border border-red-200 bg-red-50 p-4 text-red-700;
	}
</style>
