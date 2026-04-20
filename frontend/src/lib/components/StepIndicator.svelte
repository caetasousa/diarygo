<script lang="ts">
	interface Props {
		steps: string[];
		current: number; // 0-indexed
	}

	let { steps, current }: Props = $props();
</script>

<div class="steps" role="list" aria-label="Etapas do formulário">
	{#each steps as label, i}
		<div
			class="step"
			class:done={i < current}
			class:active={i === current}
			role="listitem"
			aria-current={i === current ? 'step' : undefined}
		>
			<div class="step-circle">
				{#if i < current}
					<svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
						<path d="M2 6L5 9L10 3" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
					</svg>
				{:else}
					<span class="step-num">{i + 1}</span>
				{/if}
			</div>
			<span class="step-label">{label}</span>
		</div>

		{#if i < steps.length - 1}
			<div class="step-line" class:done={i < current} aria-hidden="true"></div>
		{/if}
	{/each}
</div>

<style>
	.steps {
		display: flex;
		align-items: center;
		gap: 0;
		width: 100%;
	}

	.step {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 6px;
		flex-shrink: 0;
	}

	.step-circle {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		border: 2px solid var(--border-frost);
		background: var(--bg-card);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-tertiary);
		transition: border-color 0.15s, background 0.15s, color 0.15s;
	}

	.step.active .step-circle {
		border-color: var(--color-orange-10);
		background: rgba(255, 128, 31, 0.08);
		color: var(--color-orange-10);
	}

	.step.done .step-circle {
		border-color: var(--color-green-10);
		background: rgba(34, 197, 94, 0.08);
		color: var(--color-green-10);
	}

	.step-num {
		font-size: 0.8125rem;
		font-weight: 600;
		line-height: 1;
	}

	.step-label {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
		white-space: nowrap;
		transition: color 0.15s;
	}

	.step.active .step-label {
		color: var(--color-text-primary);
		font-weight: 500;
	}

	.step.done .step-label {
		color: var(--color-text-secondary);
	}

	.step-line {
		flex: 1;
		height: 2px;
		background: var(--border-frost);
		margin-bottom: 22px; /* alinha com o centro dos círculos */
		transition: background 0.15s;
	}

	.step-line.done {
		background: var(--color-green-4);
	}
</style>
