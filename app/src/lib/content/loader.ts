import { marked } from 'marked';
import matter from 'gray-matter';

export interface ContentMeta {
	title: string;
	slug: string;
	date: string;
	author: string;
	excerpt: string;
	tags: string[];
	featured?: boolean;
	image?: string;
	readTime?: string;
	// Whitepaper specific
	gated?: boolean;
	formId?: string;
	fileUrl?: string;
	pageCount?: number;
}

export interface Content {
	meta: ContentMeta;
	html: string;
	raw: string;
}

/**
 * Parse markdown file with frontmatter
 */
export function parseMarkdown(markdown: string): Content {
	const { data, content } = matter(markdown);
	const html = marked.parse(content) as string;

	return {
		meta: data as ContentMeta,
		html,
		raw: content
	};
}

/**
 * Load all markdown files from a glob import
 * Usage: const articles = await loadContent(import.meta.glob('/src/content/articles/*.md', { as: 'raw', eager: false }))
 */
export async function loadContent(
	modules: Record<string, () => Promise<string>>
): Promise<Content[]> {
	const content: Content[] = [];

	for (const path in modules) {
		const markdown = await modules[path]();
		const parsed = parseMarkdown(markdown);
		content.push(parsed);
	}

	// Sort by date (newest first)
	return content.sort((a, b) => new Date(b.meta.date).getTime() - new Date(a.meta.date).getTime());
}

/**
 * Load single markdown file
 */
export async function loadContentBySlug(
	modules: Record<string, () => Promise<string>>,
	slug: string
): Promise<Content | null> {
	for (const path in modules) {
		const markdown = await modules[path]();
		const parsed = parseMarkdown(markdown);
		if (parsed.meta.slug === slug) {
			return parsed;
		}
	}
	return null;
}

/**
 * Get all slugs from markdown files (for prerendering)
 */
export async function getAllSlugs(
	modules: Record<string, () => Promise<string>>
): Promise<string[]> {
	const slugs: string[] = [];

	for (const path in modules) {
		const markdown = await modules[path]();
		const { data } = matter(markdown);
		if (data.slug) {
			slugs.push(data.slug);
		}
	}

	return slugs;
}
