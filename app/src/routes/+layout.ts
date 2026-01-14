// Disable SSR globally - use client-side rendering with static fallback
export const ssr = false;

// Prerender all marketing pages for SEO
export const prerender = true;

// No trailing slashes for cleaner URLs: /projects instead of /projects/
// Creates projects.html instead of projects/index.html
export const trailingSlash = 'never';
