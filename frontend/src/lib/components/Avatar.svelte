<script lang="ts">
	interface Props {
		url?: string;
		nome: string;
		size?: 'sm' | 'md' | 'lg' | 'xl';
	}

	let { url = '', nome, size = 'md' }: Props = $props();

	let erro = $state(false);

	$effect(() => {
		url;
		erro = false;
	});

	let iniciais = $derived.by(() => {
		const partes = nome.trim().split(/\s+/).filter(Boolean);
		if (partes.length === 0) return '?';
		if (partes.length === 1) return partes[0][0].toUpperCase();
		return (partes[0][0] + partes[partes.length - 1][0]).toUpperCase();
	});

	let hue = $derived.by(() => {
		let hash = 0;
		for (let i = 0; i < nome.length; i++) hash = (hash * 31 + nome.charCodeAt(i)) | 0;
		return Math.abs(hash) % 360;
	});

	let temFoto = $derived(!!url && !erro);
</script>

<div
	class="avatar size-{size}"
	style={temFoto ? undefined : `background: hsl(${hue}, 35%, 22%); border-color: hsl(${hue}, 40%, 35%);`}
>
	{#if temFoto}
		<img
			src={url}
			alt={nome}
			loading="lazy"
			referrerpolicy="no-referrer"
			onerror={() => (erro = true)}
		/>
	{:else}
		<span class="iniciais" style="color: hsl({hue}, 80%, 75%);">{iniciais}</span>
	{/if}
</div>

<style>
	.avatar {
		border-radius: 50%;
		overflow: hidden;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid var(--border-frost);
		font-weight: 600;
		letter-spacing: -0.3px;
	}

	.size-sm { width: 32px; height: 32px; font-size: 0.75rem; }
	.size-md { width: 48px; height: 48px; font-size: 0.95rem; }
	.size-lg { width: 80px; height: 80px; font-size: 1.5rem; }
	.size-xl { width: 112px; height: 112px; font-size: 2rem; }

	.avatar img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.iniciais {
		line-height: 1;
	}
</style>
