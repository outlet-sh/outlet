#!/usr/bin/env node
/**
 * Script to format all article titles to NYT style title case
 * Usage: npx tsx scripts/formatArticleTitles.ts
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// NYT Style title case rules
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

function hasIntentionalCapitalization(word: string): boolean {
	if (word.length <= 1) return false;
	const firstChar = word.charAt(0);
	const rest = word.slice(1);
	const isStandardTitleCase = firstChar === firstChar.toUpperCase() && rest === rest.toLowerCase();
	const isAllLower = word === word.toLowerCase();
	const isAllUpper = word === word.toUpperCase();
	return !isStandardTitleCase && !isAllLower && !isAllUpper;
}

function toNYTTitleCase(title: string): string {
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

function formatArticleFile(filePath: string): void {
	const content = fs.readFileSync(filePath, 'utf-8');

	// Extract frontmatter
	const frontmatterMatch = content.match(/^---\n([\s\S]*?)\n---/);
	if (!frontmatterMatch) {
		console.log(`⚠️  No frontmatter found in ${path.basename(filePath)}`);
		return;
	}

	const frontmatter = frontmatterMatch[1];
	const titleMatch = frontmatter.match(/^title:\s*['"]?(.+?)['"]?\s*$/m);

	if (!titleMatch) {
		console.log(`⚠️  No title found in ${path.basename(filePath)}`);
		return;
	}

	const originalTitle = titleMatch[1].replace(/^['"]|['"]$/g, '');
	const formattedTitle = toNYTTitleCase(originalTitle);

	if (originalTitle === formattedTitle) {
		console.log(`✓ ${path.basename(filePath)} - already formatted`);
		return;
	}

	// Replace title in content
	const newContent = content.replace(
		/^title:\s*['"]?.*?['"]?\s*$/m,
		`title: '${formattedTitle}'`
	);

	fs.writeFileSync(filePath, newContent, 'utf-8');
	console.log(`✓ ${path.basename(filePath)}`);
	console.log(`  Old: ${originalTitle}`);
	console.log(`  New: ${formattedTitle}`);
}

// Get all article files
const articlesDir = path.join(__dirname, '../src/content/articles');
const files = fs.readdirSync(articlesDir)
	.filter(file => file.endsWith('.md'))
	.map(file => path.join(articlesDir, file));

console.log('Formatting article titles to NYT style...\n');

files.forEach(formatArticleFile);

console.log('\nDone!');
