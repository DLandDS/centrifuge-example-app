<script lang="ts">
	import { onMount } from 'svelte';
	import { authStore, initializeAuth } from '$lib/stores/auth';
	import Login from '$lib/components/Login.svelte';
	import Chat from '$lib/components/Chat.svelte';

	onMount(() => {
		initializeAuth();
	});
</script>

<svelte:head>
	<title>Centrifuge Chat App</title>
	<meta name="description" content="Real-time chat application using Centrifuge" />
</svelte:head>

{#if $authStore.isAuthenticated}
	<Chat />
{:else}
	<Login />
{/if}

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
		background-color: #f5f5f5;
	}

	:global(*) {
		box-sizing: border-box;
	}
</style>
