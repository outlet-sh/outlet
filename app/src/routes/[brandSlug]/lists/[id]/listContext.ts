import { getContext, setContext } from 'svelte';
import type { ListInfo } from '$lib/api';

const LIST_CONTEXT_KEY = Symbol('list-context');

export interface ListContext {
	readonly list: ListInfo | null;
	readonly listId: string;
	readonly basePath: string;
	reload: () => Promise<void>;
}

export function setListContext(ctx: ListContext) {
	setContext(LIST_CONTEXT_KEY, ctx);
}

export function getListContext(): ListContext {
	const ctx = getContext<ListContext>(LIST_CONTEXT_KEY);
	if (!ctx) {
		throw new Error('List context not found. Make sure this component is rendered under the list layout.');
	}
	return ctx;
}
