<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { authStore, logout } from '$lib/stores/auth';
	import { chatStore, connectToCentrifuge, subscribeToTopic, sendMessage, disconnect } from '$lib/stores/chat';

	let messageContent = '';
	let messagesContainer: HTMLElement;

	$: if ($chatStore.messages.length > 0 && messagesContainer) {
		setTimeout(() => messagesContainer.scrollTop = messagesContainer.scrollHeight, 100);
	}

	onMount(() => {
		if ($authStore.centrifugeToken) {
			connectToCentrifuge($authStore.centrifugeToken);
			// Subscribe to default topic
			subscribeToTopic('general');
		}
	});

	onDestroy(() => {
		disconnect();
	});

	const handleTopicChange = (topic: string) => {
		subscribeToTopic(topic);
	};

	const handleSendMessage = () => {
		if (messageContent.trim() && $authStore.user) {
			sendMessage(messageContent.trim(), $authStore.user.username);
			messageContent = '';
		}
	};

	const handleLogout = () => {
		disconnect();
		logout();
	};

	const formatTime = (timestamp: number) => {
		return new Date(timestamp).toLocaleTimeString();
	};
</script>

<div class="chat-container">
	<header class="chat-header">
		<div class="user-info">
			<span>Welcome, {$authStore.user?.username}!</span>
		</div>
		<div class="connection-status" class:connected={$chatStore.isConnected}>
			{$chatStore.isConnected ? 'Connected' : 'Disconnected'}
		</div>
		<button class="logout-btn" on:click={handleLogout}>Logout</button>
	</header>

	<div class="chat-content">
		<aside class="topics-sidebar">
			<h3>Topics</h3>
			<ul class="topics-list">
				{#each $chatStore.topics as topic}
					<li>
						<button
							class="topic-btn"
							class:active={$chatStore.activeTopic === topic}
							on:click={() => handleTopicChange(topic)}
						>
							#{topic}
						</button>
					</li>
				{/each}
			</ul>
		</aside>

		<main class="chat-main">
			<div class="current-topic">
				<h2>#{$chatStore.activeTopic || 'Select a topic'}</h2>
			</div>

			<div class="messages" bind:this={messagesContainer}>
				{#each $chatStore.messages as message (message.id)}
					<div class="message">
						<div class="message-header">
							<span class="username">{message.username}</span>
							<span class="timestamp">{formatTime(message.timestamp)}</span>
						</div>
						<div class="message-content">{message.content}</div>
					</div>
				{/each}
				{#if $chatStore.messages.length === 0}
					<div class="no-messages">
						{#if $chatStore.activeTopic}
							No messages in #{$chatStore.activeTopic} yet. Start the conversation!
						{:else}
							Select a topic to start chatting
						{/if}
					</div>
				{/if}
			</div>

			{#if $chatStore.activeTopic}
				<div class="message-input">
					<form on:submit|preventDefault={handleSendMessage}>
						<input
							type="text"
							bind:value={messageContent}
							placeholder="Type your message..."
							disabled={!$chatStore.isConnected}
						/>
						<button
							type="submit"
							disabled={!$chatStore.isConnected || !messageContent.trim()}
						>
							Send
						</button>
					</form>
				</div>
			{/if}
		</main>
	</div>
</div>

<style>
	.chat-container {
		height: 100vh;
		display: flex;
		flex-direction: column;
	}

	.chat-header {
		background-color: #333;
		color: white;
		padding: 1rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.connection-status {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		background-color: #666;
		font-size: 0.9rem;
	}

	.connection-status.connected {
		background-color: #4caf50;
	}

	.logout-btn {
		background-color: #f44336;
		color: white;
		border: none;
		padding: 0.5rem 1rem;
		border-radius: 4px;
		cursor: pointer;
	}

	.logout-btn:hover {
		background-color: #d32f2f;
	}

	.chat-content {
		flex: 1;
		display: flex;
		overflow: hidden;
	}

	.topics-sidebar {
		width: 200px;
		background-color: #f5f5f5;
		border-right: 1px solid #ddd;
		padding: 1rem;
	}

	.topics-sidebar h3 {
		margin: 0 0 1rem 0;
		color: #333;
	}

	.topics-list {
		list-style: none;
		padding: 0;
		margin: 0;
	}

	.topics-list li {
		margin-bottom: 0.5rem;
	}

	.topic-btn {
		width: 100%;
		text-align: left;
		background: none;
		border: none;
		padding: 0.5rem;
		border-radius: 4px;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.topic-btn:hover {
		background-color: #e0e0e0;
	}

	.topic-btn.active {
		background-color: #0066cc;
		color: white;
	}

	.chat-main {
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.current-topic {
		padding: 1rem;
		border-bottom: 1px solid #ddd;
		background-color: white;
	}

	.current-topic h2 {
		margin: 0;
		color: #333;
	}

	.messages {
		flex: 1;
		overflow-y: auto;
		padding: 1rem;
		background-color: #fafafa;
	}

	.message {
		background-color: white;
		border-radius: 8px;
		padding: 0.75rem;
		margin-bottom: 0.5rem;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
	}

	.message-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.5rem;
	}

	.username {
		font-weight: bold;
		color: #0066cc;
	}

	.timestamp {
		font-size: 0.8rem;
		color: #666;
	}

	.message-content {
		color: #333;
		line-height: 1.4;
	}

	.no-messages {
		text-align: center;
		color: #666;
		font-style: italic;
		margin-top: 2rem;
	}

	.message-input {
		padding: 1rem;
		background-color: white;
		border-top: 1px solid #ddd;
	}

	.message-input form {
		display: flex;
		gap: 0.5rem;
	}

	.message-input input {
		flex: 1;
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
	}

	.message-input input:focus {
		outline: none;
		border-color: #0066cc;
	}

	.message-input button {
		background-color: #0066cc;
		color: white;
		border: none;
		padding: 0.75rem 1.5rem;
		border-radius: 4px;
		cursor: pointer;
		font-size: 1rem;
	}

	.message-input button:hover:not(:disabled) {
		background-color: #0052a3;
	}

	.message-input button:disabled {
		background-color: #ccc;
		cursor: not-allowed;
	}
</style>