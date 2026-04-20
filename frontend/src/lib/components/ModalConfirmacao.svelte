<script lang="ts">
	interface Props {
		open: boolean;
		titulo: string;
		mensagem: string;
		textoConfirmar?: string;
		tipo?: 'danger' | 'default';
		onConfirm: () => void;
		onCancel: () => void;
	}

	let {
		open = $bindable(false),
		titulo,
		mensagem,
		textoConfirmar = 'Confirmar',
		tipo = 'default',
		onConfirm,
		onCancel
	}: Props = $props();

	function onKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onCancel();
	}
</script>

<svelte:window onkeydown={onKeydown} />

{#if open}
	<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
	<div class="overlay" onclick={onCancel} role="presentation">
		<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
		<!-- svelte-ignore a11y_interactive_supports_focus -->
		<div
			class="modal"
			role="dialog"
			aria-modal="true"
			aria-labelledby="modal-titulo"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 id="modal-titulo" class="modal-titulo">{titulo}</h2>
			<p class="modal-mensagem">{mensagem}</p>
			<div class="modal-acoes">
				<!-- svelte-ignore a11y_autofocus -->
				<button type="button" class="btn-cancelar" onclick={onCancel} autofocus>
					Cancelar
				</button>
				<button
					type="button"
					class="btn-confirmar"
					class:danger={tipo === 'danger'}
					onclick={onConfirm}
				>
					{textoConfirmar}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		backdrop-filter: blur(2px);
	}

	.modal {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-large);
		padding: 28px 24px 24px;
		max-width: 420px;
		width: calc(100% - 32px);
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.modal-titulo {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
		color: var(--color-text-primary);
		letter-spacing: -0.2px;
	}

	.modal-mensagem {
		margin: 0;
		font-size: 0.875rem;
		color: var(--color-text-secondary);
		line-height: 1.5;
	}

	.modal-acoes {
		display: flex;
		gap: 10px;
		justify-content: flex-end;
		margin-top: 8px;
	}

	.btn-cancelar,
	.btn-confirmar {
		padding: 9px 18px;
		border-radius: var(--radius-pill);
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.12s, border-color 0.12s;
	}

	.btn-cancelar {
		background: transparent;
		border: 1px solid var(--border-frost);
		color: var(--color-text-primary);
	}
	.btn-cancelar:hover {
		background: var(--bg-hover-subtle);
		border-color: var(--border-frost-hover);
	}

	.btn-confirmar {
		background: var(--color-orange-10);
		border: 1px solid var(--color-orange-10);
		color: #fff;
	}
	.btn-confirmar:hover { background: var(--color-orange-11); border-color: var(--color-orange-11); }
	.btn-confirmar.danger {
		background: #ef4444;
		border-color: #ef4444;
	}
	.btn-confirmar.danger:hover { background: #dc2626; border-color: #dc2626; }
</style>
