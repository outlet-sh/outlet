/**
 * Convert a title to NYT style title case
 *
 * NYT Style Rules:
 * - Capitalize major words (nouns, pronouns, verbs, adjectives, adverbs, subordinate conjunctions)
 * - Capitalize the first and last word
 * - Lowercase articles (a, an, the)
 * - Lowercase coordinating conjunctions (and, but, or, nor, for, yet, so)
 * - Lowercase prepositions (at, by, for, in, of, on, to, up, via, etc.)
 */

const articles = ['a', 'an', 'the'];
const coordinatingConjunctions = ['and', 'but', 'or', 'nor', 'for', 'yet', 'so'];
const prepositions = [
	'at', 'by', 'for', 'in', 'of', 'on', 'to', 'up', 'via',
	'with', 'from', 'into', 'upon', 'over', 'after', 'around',
	'before', 'behind', 'below', 'beneath', 'beside', 'between',
	'beyond', 'during', 'inside', 'near', 'off', 'out', 'outside',
	'through', 'throughout', 'toward', 'under', 'underneath', 'until',
	'within', 'without', 'about', 'above', 'across', 'against', 'along',
	'among', 'amongst', 'anti', 'aside'
];

const lowercaseWords = [...articles, ...coordinatingConjunctions, ...prepositions];

// Common acronyms to preserve in uppercase
const acronyms = ['AI', 'MCP', 'API', 'SDK', 'URL', 'HTTP', 'HTTPS', 'REST', 'SQL', 'JSON', 'XML', 'HTML', 'CSS', 'JS', 'TS', 'UI', 'UX', 'AWS', 'GCP', 'RAG', 'LLM', 'GPT', 'GoMCP', 'NYT', 'APA', 'MLA'];

export function toNYTTitleCase(title: string): string {
	const words = title.trim().split(/\s+/);

	return words.map((word, index) => {
		const cleanWord = word.replace(/[^\w]/g, '');

		// Preserve intentional capitalization (mixed case words like GoMCP, iPhone, etc.)
		if (hasIntentionalCapitalization(cleanWord)) {
			return word;
		}

		// Check for known acronyms
		if (acronyms.includes(cleanWord.toUpperCase())) {
			return cleanWord.toUpperCase();
		}

		// Always capitalize first and last word
		if (index === 0 || index === words.length - 1) {
			return capitalizeWord(word);
		}

		// Check if word should be lowercase
		if (lowercaseWords.includes(cleanWord.toLowerCase())) {
			return word.toLowerCase();
		}

		// Capitalize everything else
		return capitalizeWord(word);
	}).join(' ');
}

function hasIntentionalCapitalization(word: string): boolean {
	// If word has mixed case (not all lowercase, not all uppercase, not just first letter capitalized)
	// then it's likely intentional (e.g., GoMCP, iPhone, eBay)
	if (word.length <= 1) return false;

	const firstChar = word.charAt(0);
	const rest = word.slice(1);

	// Check if it's standard title case (First letter capital, rest lowercase)
	const isStandardTitleCase = firstChar === firstChar.toUpperCase() && rest === rest.toLowerCase();

	// Check if it's all lowercase or all uppercase
	const isAllLower = word === word.toLowerCase();
	const isAllUpper = word === word.toUpperCase();

	// If it's not standard title case and not all one case, it's intentional
	return !isStandardTitleCase && !isAllLower && !isAllUpper;
}

function capitalizeWord(word: string): string {
	if (!word) return word;

	// Handle hyphenated words - capitalize both parts
	if (word.includes('-')) {
		return word.split('-')
			.map(part => part.charAt(0).toUpperCase() + part.slice(1).toLowerCase())
			.join('-');
	}

	// Handle words with apostrophes (don't, it's, etc.)
	if (word.includes("'")) {
		const parts = word.split("'");
		return parts[0].charAt(0).toUpperCase() + parts[0].slice(1).toLowerCase() + "'" + parts.slice(1).join("'").toLowerCase();
	}

	// Standard capitalization
	return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
}

/**
 * Format article frontmatter title to NYT style
 */
export function formatArticleTitle(title: string): string {
	return toNYTTitleCase(title);
}
