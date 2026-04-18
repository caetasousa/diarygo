<script lang="ts">
	interface Props {
		url: string;
		alt: string;
		size?: 'sm' | 'md' | 'lg';
		rounded?: 'circle' | 'square';
		clickable?: boolean;
		onClick?: () => void;
	}

	let {
		url,
		alt,
		size = 'md',
		rounded = 'square',
		clickable = false,
		onClick
	}: Props = $props();

	let erro = $state(false);
	let carregando = $state(true);

	$effect(() => {
		url;
		erro = false;
		carregando = !!url;
	});

	function onLoad() {
		carregando = false;
	}

	function onError() {
		erro = true;
		carregando = false;
	}
</script>

{#if clickable}
	<button
		type="button"
		class="preview size-{size} shape-{rounded} clickable"
		onclick={onClick}
	>
		{#if !url || erro}
			<div class="fallback" aria-label={alt}>
				<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
					<rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
					<circle cx="8.5" cy="8.5" r="1.5"/>
					<polyline points="21 15 16 10 5 21"/>
				</svg>
			</div>
		{:else}
			{#if carregando}
				<div class="skeleton"></div>
			{/if}
			<img
				src={url}
				{alt}
				loading="lazy"
				referrerpolicy="no-referrer"
				onload={onLoad}
				onerror={onError}
				class:loaded={!carregando}
			/>
		{/if}
	</button>
{:else}
<div class="preview size-{size} shape-{rounded}">
	{#if !url || erro}
		<div class="fallback" aria-label={alt}>
			<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
				<rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
				<circle cx="8.5" cy="8.5" r="1.5"/>
				<polyline points="21 15 16 10 5 21"/>
			</svg>
		</div>
	{:else}
		{#if carregando}
			<div class="skeleton"></div>
		{/if}
		<img
			src={url}
			{alt}
			loading="lazy"
			referrerpolicy="no-referrer"
			onload={onLoad}
			onerror={onError}
			class:loaded={!carregando}
		/>
	{/if}
</div>
{/if}

<style>
	.preview {
		position: relative;
		overflow: hidden;
		background: rgba(255, 255, 255, 0.04);
		border: 1px solid var(--border-frost);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-tertiary);
		padding: 0;
	}
	button.preview { font: inherit; }

	.size-sm { width: 40px; height: 40px; }
	.size-md { width: 88px; height: 88px; }
	.size-lg { width: 140px; height: 140px; }

	.shape-square { border-radius: 10px; }
	.shape-circle { border-radius: 50%; }

	.clickable { cursor: pointer; transition: transform 0.12s, border-color 0.12s; }
	.clickable:hover { transform: translateY(-1px); border-color: var(--color-blue-5); }
	.clickable:focus-visible { outline: 2px solid var(--color-blue-10); outline-offset: 2px; }

	.preview img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		opacity: 0;
		transition: opacity 0.2s ease;
	}
	.preview img.loaded { opacity: 1; }

	.skeleton {
		position: absolute;
		inset: 0;
		background: linear-gradient(90deg, rgba(255,255,255,0.03), rgba(255,255,255,0.08), rgba(255,255,255,0.03));
		background-size: 200% 100%;
		animation: shimmer 1.4s infinite;
	}
	@keyframes shimmer { from { background-position: 200% 0; } to { background-position: -200% 0; } }

	.fallback {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
	}
</style>
