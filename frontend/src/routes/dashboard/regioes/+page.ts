import { redirect } from '@sveltejs/kit';
import { browser } from '$app/environment';
import { get } from 'svelte/store';
import { isAuthenticated, auth } from '$lib/stores/auth';

export function load() {
	if (browser) {
		if (!get(isAuthenticated)) throw redirect(302, '/login');
		const tipo = get(auth)?.payload?.tipo;
		if (tipo !== 'PROFISSIONAL') throw redirect(302, '/dashboard');
	}
}
