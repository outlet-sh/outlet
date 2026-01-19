import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ params }) => {
	return {
		listId: params.id,
		orgSlug: params.orgSlug
	};
};
