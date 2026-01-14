/**
 * SEO utilities for generating meta tags
 */

export interface SEOConfig {
	title?: string;
	description?: string;
	image?: string;
	url?: string;
	type?: string;
	siteName?: string;
}

const defaultConfig: SEOConfig = {
	title: 'Your Platform',
	description: 'Self-hosted SaaS platform for building amazing products.',
	siteName: 'Your Platform',
	type: 'website'
};

export function setSEO(config: SEOConfig = {}): SEOConfig {
	return {
		...defaultConfig,
		...config,
		title: config.title ? `${config.title} | Your Platform` : defaultConfig.title
	};
}

export function generateMetaTags(config: SEOConfig) {
	const seo = setSEO(config);
	return {
		title: seo.title,
		description: seo.description,
		'og:title': seo.title,
		'og:description': seo.description,
		'og:type': seo.type,
		'og:site_name': seo.siteName,
		'og:image': seo.image,
		'og:url': seo.url,
		'twitter:card': 'summary_large_image',
		'twitter:title': seo.title,
		'twitter:description': seo.description,
		'twitter:image': seo.image
	};
}
