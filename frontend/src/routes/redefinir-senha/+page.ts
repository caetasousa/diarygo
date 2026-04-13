import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ url }) => {
	const token = url.searchParams.get('token');
	if (!token) {
		throw redirect(302, '/recuperar-senha');
	}
	return { token };
};
