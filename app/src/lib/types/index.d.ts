// Global type definitions

// Google Ads gtag
declare function gtag(...args: any[]): void;

export interface ContactFormData {
	name: string;
	email: string;
	company?: string;
	employees?: string;
	message?: string;
}

export interface InsightPost {
	slug: string;
	title: string;
	description: string;
	publishedAt: string;
	content?: string;
}

export interface CaseStudy {
	id: string;
	company: string;
	description: string;
	results: string[];
}

export interface ServiceArea {
	title: string;
	description: string;
	icon?: string;
}

// Export all call-related types
export * from './call';
