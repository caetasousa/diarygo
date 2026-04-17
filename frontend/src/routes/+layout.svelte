<script lang="ts">
	import '../app.css';
	import { page } from '$app/state';
	import Navbar from '$lib/components/Navbar.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import ToastContainer from '$lib/components/ToastContainer.svelte';

	let { children } = $props();

	let isDashboard = $derived(page.url.pathname.startsWith('/dashboard'));
</script>

{#if !isDashboard}
	<Navbar />
{/if}

<main class:dashboard={isDashboard}>
	{@render children()}
</main>

{#if !isDashboard}
	<Footer />
{/if}

<ToastContainer />

<style>
	main {
		min-height: calc(100vh - 60px);
	}

	main.dashboard {
		min-height: 100vh;
	}
</style>
