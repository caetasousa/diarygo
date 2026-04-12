import { writable, derived } from 'svelte/store';
import type { TokenPayload, AuthState } from '$lib/types';
import { api } from '$lib/api/client';
import { browser } from '$app/environment';

const TOKEN_KEY = 'diarygo_token';

function parsePayload(token: string): TokenPayload | null {
	try {
		const [, payload] = token.split('.');
		return JSON.parse(atob(payload)) as TokenPayload;
	} catch {
		return null;
	}
}

function isExpired(payload: TokenPayload): boolean {
	return Date.now() / 1000 > payload.exp;
}

function createAuthStore() {
	const initialToken = browser ? localStorage.getItem(TOKEN_KEY) : null;
	const initialPayload = initialToken ? parsePayload(initialToken) : null;
	const validToken =
		initialToken && initialPayload && !isExpired(initialPayload) ? initialToken : null;

	const { subscribe, set, update } = writable<AuthState>({
		token: validToken,
		payload: validToken ? initialPayload : null,
		loading: false
	});

	if (validToken) {
		api.setToken(validToken);
	} else if (browser) {
		localStorage.removeItem(TOKEN_KEY);
	}

	return {
		subscribe,

		login(token: string) {
			const payload = parsePayload(token);
			if (!payload || isExpired(payload)) return;
			if (browser) localStorage.setItem(TOKEN_KEY, token);
			api.setToken(token);
			set({ token, payload, loading: false });
		},

		logout() {
			if (browser) localStorage.removeItem(TOKEN_KEY);
			api.setToken(null);
			set({ token: null, payload: null, loading: false });
		},

		setLoading(loading: boolean) {
			update((s) => ({ ...s, loading }));
		}
	};
}

export const auth = createAuthStore();

export const isAuthenticated = derived(auth, ($auth) => !!$auth.token);
export const currentUser = derived(auth, ($auth) => $auth.payload);
