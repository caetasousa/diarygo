import { redirect } from '@sveltejs/kit';
import { browser } from '$app/environment';
import { get } from 'svelte/store';
import { isAuthenticated, currentUser } from '$lib/stores/auth';

export function load() {
	if (browser) {
		if (!get(isAuthenticated)) throw redirect(302, '/login');
		const user = get(currentUser);
		if (user?.tipo !== 'CLIENTE') throw redirect(302, '/dashboard');
	}
}
