<script lang="ts">
	import '../app.css';
	import Navbar from '$lib/components/Navbar.svelte';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
</script>

<Navbar />

<main>
	<slot />
</main>

<ToastContainer />

<style>
	main {
		min-height: calc(100vh - 60px);
	}
</style>
