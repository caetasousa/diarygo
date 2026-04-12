<script lang="ts">
	import Toast from './Toast.svelte';
	import { toasts } from '$lib/stores/toasts';
</script>

<div class="toast-container">
	{#each $toasts as toast (toast.id)}
		<Toast
			message={toast.message}
			type={toast.type}
			on:close={() => toasts.remove(toast.id)}
		/>
	{/each}
</div>

<style>
	.toast-container {
		position: fixed;
		bottom: var(--space-6);
		right: var(--space-6);
		z-index: 1000;
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
		align-items: flex-end;
	}

	@media (max-width: 480px) {
		.toast-container {
			left: var(--space-4);
			right: var(--space-4);
			bottom: var(--space-4);
			align-items: stretch;
		}
	}
</style>
