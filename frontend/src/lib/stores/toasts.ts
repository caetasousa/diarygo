import { writable } from 'svelte/store';

export interface ToastItem {
	id: string;
	message: string;
	type: 'success' | 'error' | 'info';
}

function createToasts() {
	const { subscribe, update } = writable<ToastItem[]>([]);

	function add(message: string, type: ToastItem['type'] = 'info') {
		const id = crypto.randomUUID();
		update((list) => [...list, { id, message, type }]);
		return id;
	}

	function remove(id: string) {
		update((list) => list.filter((t) => t.id !== id));
	}

	return {
		subscribe,
		success: (msg: string) => add(msg, 'success'),
		error: (msg: string) => add(msg, 'error'),
		info: (msg: string) => add(msg, 'info'),
		remove
	};
}

export const toasts = createToasts();
