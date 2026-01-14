/**
 * Funnel Analytics - Page-flow tracking for conversion optimization
 * Tracks: Entry -> Apply -> Calendar -> Completion
 */

import { browser } from '$app/environment';
import { trackEvent } from '$lib/api';
import type { EventRequest } from '$lib/api';

const VISITOR_ID_KEY = 'funnel_visitor_id';
const SESSION_ID_KEY = 'funnel_session_id';

export type FunnelStep =
	| 'entry'
	| 'apply'
	| 'calendar'
	| 'completion'
	| 'thank_you_view'
	| 'confirm_pending'
	| 'optin_complete'
	| 'workshop_checkout_view'
	| 'product_page_view'
	| 'quiz_started'
	| 'quiz_completed'
	| 'upsell_view'
	| 'downsell_view'
	| 'order_bump_view'
	| 'upsell_accepted'
	| 'downsell_accepted'
	| 'order_bump_accepted'
	| 'upsell_declined'
	| 'downsell_declined'
	| 'order_bump_declined';

/**
 * Configuration for a funnel step page
 * Export this from +page.ts to define funnel metadata
 */
export interface FunnelStepConfig {
	/** Position in funnel sequence (1 = entry point) */
	position: number;
	/** Step type for categorization */
	type: 'opt-in' | 'confirm' | 'thank-you' | 'checkout' | 'upsell' | 'downsell' | 'quiz' | 'workshop';
	/** Display title for admin UI */
	title: string;
	/** Short description */
	description?: string;
	/** Access requirement - what step must be completed first */
	requires?: 'opt-in' | 'email_verified' | 'purchase' | string;
	/** Path to downloadable lead magnet (if any) */
	leadMagnet?: string;
	/** Human-readable lead magnet name */
	leadMagnetName?: string;
}

// Generate a simple UUID
function generateId(): string {
	return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
		const r = (Math.random() * 16) | 0;
		const v = c === 'x' ? r : (r & 0x3) | 0x8;
		return v.toString(16);
	});
}

// Get or create visitor ID (persisted in localStorage)
function getVisitorId(): string {
	if (!browser) return '';

	let visitorId = localStorage.getItem(VISITOR_ID_KEY);
	if (!visitorId) {
		visitorId = generateId();
		localStorage.setItem(VISITOR_ID_KEY, visitorId);
	}
	return visitorId;
}

// Get or create session ID (persisted in sessionStorage)
function getSessionId(): string {
	if (!browser) return '';

	let sessionId = sessionStorage.getItem(SESSION_ID_KEY);
	if (!sessionId) {
		sessionId = generateId();
		sessionStorage.setItem(SESSION_ID_KEY, sessionId);
	}
	return sessionId;
}

// Parse UTM parameters from URL
function getUtmParams(): { source?: string; medium?: string; campaign?: string } {
	if (!browser) return {};

	const params = new URLSearchParams(window.location.search);
	return {
		source: params.get('utm_source') || undefined,
		medium: params.get('utm_medium') || undefined,
		campaign: params.get('utm_campaign') || undefined
	};
}

// Send funnel event to backend
async function sendFunnelEvent(data: EventRequest): Promise<void> {
	try {
		await trackEvent(data);
	} catch (error) {
		console.error('Failed to send funnel event:', error);
	}
}

// Track scroll state per page
const scrollTrackers = new Map<string, { maxDepth: number; lastSent: number }>();

/**
 * Initialize funnel tracking - call on app mount
 */
export function initFunnel(): { visitorId: string; sessionId: string } {
	return {
		visitorId: getVisitorId(),
		sessionId: getSessionId()
	};
}

/**
 * Track a funnel step
 */
export function trackFunnelStep(step: FunnelStep, page?: string, leadId?: string): void {
	if (!browser) return;

	const utm = getUtmParams();
	const data: EventRequest = {
		event: `funnel_${step}`,
		visitor_id: getVisitorId(),
		session_id: getSessionId(),
		funnel_step: step,
		page: page || window.location.pathname,
		referrer: document.referrer || undefined,
		properties: {
			utm_source: utm.source,
			utm_medium: utm.medium,
			utm_campaign: utm.campaign,
			lead_id: leadId
		}
	};

	sendFunnelEvent(data);
}

/**
 * Start scroll depth tracking for a page
 * Returns cleanup function
 */
export function trackScrollDepth(step: FunnelStep = 'entry', page?: string): () => void {
	if (!browser) return () => {};

	const pagePath = page || window.location.pathname;

	// Initialize tracker for this page
	if (!scrollTrackers.has(pagePath)) {
		scrollTrackers.set(pagePath, { maxDepth: 0, lastSent: 0 });
	}

	const tracker = scrollTrackers.get(pagePath)!;

	const calculateScrollDepth = (): number => {
		const scrollTop = window.scrollY;
		const docHeight = document.documentElement.scrollHeight - window.innerHeight;
		if (docHeight <= 0) return 100;
		return Math.min(100, Math.round((scrollTop / docHeight) * 100));
	};

	const handleScroll = (): void => {
		const depth = calculateScrollDepth();

		// Only update if we've scrolled deeper
		if (depth > tracker.maxDepth) {
			tracker.maxDepth = depth;

			// Send updates at milestones (25%, 50%, 75%, 100%) or if 10% increase since last send
			const milestones = [25, 50, 75, 100];
			const isAtMilestone = milestones.includes(depth);
			const significantIncrease = depth - tracker.lastSent >= 10;

			if (isAtMilestone || significantIncrease) {
				tracker.lastSent = depth;

				const utm = getUtmParams();
				sendFunnelEvent({
					event: 'scroll_depth',
					visitor_id: getVisitorId(),
					session_id: getSessionId(),
					funnel_step: step,
					page: pagePath,
					scroll_depth: depth,
					properties: {
						utm_source: utm.source,
						utm_medium: utm.medium,
						utm_campaign: utm.campaign
					}
				});
			}
		}
	};

	// Throttled scroll handler
	let ticking = false;
	const throttledScroll = (): void => {
		if (!ticking) {
			requestAnimationFrame(() => {
				handleScroll();
				ticking = false;
			});
			ticking = true;
		}
	};

	window.addEventListener('scroll', throttledScroll, { passive: true });

	// Return cleanup function
	return () => {
		window.removeEventListener('scroll', throttledScroll);
		scrollTrackers.delete(pagePath);
	};
}

/**
 * Get the current funnel session info
 */
export function getFunnelSession(): { visitorId: string; sessionId: string } | null {
	if (!browser) return null;

	const visitorId = localStorage.getItem(VISITOR_ID_KEY);
	const sessionId = sessionStorage.getItem(SESSION_ID_KEY);

	if (!visitorId || !sessionId) return null;
	return { visitorId, sessionId };
}
