<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let message: string;
	export let type: 'success' | 'error' | 'info' = 'info';
	export let duration = 4000;

	const dispatch = createEventDispatcher();

	let visible = true;

	$: {
		if (duration > 0) {
			setTimeout(() => {
				visible = false;
				setTimeout(() => dispatch('close'), 300);
			}, duration);
		}
	}
</script>

{#if visible}
	<div class="toast toast-{type} animate-in" role="alert">
		<span class="toast-icon">
			{#if type === 'success'}✓{:else if type === 'error'}✕{:else}ℹ{/if}
		</span>
		<span class="toast-message">{message}</span>
		<button class="toast-close" on:click={() => { visible = false; dispatch('close'); }}>✕</button>
	</div>
{/if}

<style>
	.toast {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: 12px 16px;
		border-radius: var(--radius-standard);
		border: 1px solid var(--border-frost);
		background: rgba(10, 10, 10, 0.95);
		backdrop-filter: blur(8px);
		font-size: 0.875rem;
		min-width: 280px;
		max-width: 400px;
		box-shadow: var(--shadow-ring);
	}

	.toast-success .toast-icon { color: var(--color-green-10); }
	.toast-error .toast-icon   { color: var(--color-red-10); }
	.toast-info .toast-icon    { color: var(--color-blue-10); }

	.toast-message {
		flex: 1;
		color: var(--color-text-primary);
	}

	.toast-close {
		background: none;
		border: none;
		color: var(--color-text-tertiary);
		cursor: pointer;
		font-size: 0.75rem;
		padding: 2px;
		line-height: 1;
		transition: color 0.15s ease;
	}

	.toast-close:hover {
		color: var(--color-text-primary);
	}
</style>
