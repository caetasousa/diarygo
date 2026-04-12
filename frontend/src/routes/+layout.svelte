<script lang="ts">
	import '../app.css';
	import Navbar from '$lib/components/Navbar.svelte';
	import ToastContainer from '$lib/components/ToastContainer.svelte';

	let { children } = $props();
</script>

<Navbar />

<main>
	{@render children()}
</main>

<ToastContainer />

<style>
	main {
		min-height: calc(100vh - 60px);
	}
</style>
