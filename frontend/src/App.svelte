<script>
	import { onMount } from 'svelte';
	import { Centrifuge } from 'centrifuge';

	let username = '';
	let password = '';
	let isLoggedIn = false;
	let token = '';
	let centrifuge = null;
	let currentTopic = 'general';
	let messageContent = '';
	let messages = [];
	let subscribedTopics = new Set();
	let availableTopics = ['general', 'tech', 'random', 'announcements'];

	const API_BASE = 'http://localhost:8080';
	const WS_URL = 'ws://localhost:8080/connection/websocket';

	async function login() {
		try {
			const response = await fetch(`${API_BASE}/api/login`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({ username, password }),
			});

			if (response.ok) {
				const data = await response.json();
				token = data.token;
				isLoggedIn = true;
				connectToCentrifuge();
			} else {
				alert('Login failed');
			}
		} catch (error) {
			console.error('Login error:', error);
			alert('Login error');
		}
	}

	function connectToCentrifuge() {
		centrifuge = new Centrifuge(WS_URL, {
			token: token,
			debug: true
		});

		centrifuge.on('connected', function(ctx) {
			console.log('Connected to Centrifuge:', ctx);
		});

		centrifuge.on('disconnected', function(ctx) {
			console.log('Disconnected from Centrifuge:', ctx);
		});

		centrifuge.on('error', function(ctx) {
			console.error('Centrifuge error:', ctx);
		});

		centrifuge.connect();

		// Subscribe to the current topic
		subscribeToTopic(currentTopic);
	}

	function subscribeToTopic(topic) {
		if (!centrifuge || subscribedTopics.has(topic)) {
			return;
		}

		const subscription = centrifuge.newSubscription(topic);

		subscription.on('subscribed', function(ctx) {
			console.log('Subscribed to topic:', topic);
			subscribedTopics.add(topic);
		});

		subscription.on('publication', function(ctx) {
			const message = ctx.data;
			messages = [...messages, message];
		});

		subscription.subscribe();
	}

	function switchTopic(topic) {
		currentTopic = topic;
		messages = []; // Clear messages when switching topics
		subscribeToTopic(topic);
	}

	async function sendMessage() {
		if (!messageContent.trim() || !isLoggedIn) {
			return;
		}

		try {
			const response = await fetch(`${API_BASE}/api/topics/${currentTopic}/messages`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${token}`,
				},
				body: JSON.stringify({ content: messageContent }),
			});

			if (response.ok) {
				messageContent = '';
			} else {
				alert('Failed to send message');
			}
		} catch (error) {
			console.error('Send message error:', error);
			alert('Failed to send message');
		}
	}

	function logout() {
		if (centrifuge) {
			centrifuge.disconnect();
		}
		isLoggedIn = false;
		token = '';
		messages = [];
		subscribedTopics.clear();
	}

	function formatTime(timestamp) {
		return new Date(timestamp).toLocaleTimeString();
	}
</script>

<main>
	<div class="container">
		<h1>🚀 Centrifuge Chat App</h1>
		
		{#if !isLoggedIn}
			<div class="login-form">
				<h2>Login</h2>
				<form on:submit|preventDefault={login}>
					<input
						type="text"
						placeholder="Username"
						bind:value={username}
						required
					/>
					<input
						type="password"
						placeholder="Password (optional)"
						bind:value={password}
					/>
					<button type="submit">Login</button>
				</form>
			</div>
		{:else}
			<div class="chat-container">
				<div class="header">
					<span>Welcome, {username}!</span>
					<button on:click={logout} class="logout-btn">Logout</button>
				</div>
				
				<div class="topics">
					<h3>Topics</h3>
					<div class="topic-list">
						{#each availableTopics as topic}
							<button
								class="topic-btn {currentTopic === topic ? 'active' : ''}"
								on:click={() => switchTopic(topic)}
							>
								#{topic}
							</button>
						{/each}
					</div>
				</div>

				<div class="messages-container">
					<div class="messages-header">
						<h3>#{currentTopic}</h3>
					</div>
					
					<div class="messages">
						{#each messages as message}
							<div class="message">
								<div class="message-header">
									<span class="username">{message.username}</span>
									<span class="timestamp">{formatTime(message.timestamp)}</span>
								</div>
								<div class="message-content">{message.content}</div>
							</div>
						{/each}
					</div>

					<div class="message-input">
						<form on:submit|preventDefault={sendMessage}>
							<input
								type="text"
								placeholder="Type your message..."
								bind:value={messageContent}
								required
							/>
							<button type="submit">Send</button>
						</form>
					</div>
				</div>
			</div>
		{/if}
	</div>
</main>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
		background-color: #f5f5f5;
	}

	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 20px;
	}

	h1 {
		text-align: center;
		color: #333;
		margin-bottom: 30px;
	}

	.login-form {
		max-width: 400px;
		margin: 0 auto;
		background: white;
		padding: 30px;
		border-radius: 8px;
		box-shadow: 0 2px 10px rgba(0,0,0,0.1);
	}

	.login-form h2 {
		text-align: center;
		margin-bottom: 20px;
		color: #333;
	}

	.login-form form {
		display: flex;
		flex-direction: column;
		gap: 15px;
	}

	.login-form input {
		padding: 12px;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 16px;
	}

	.login-form button {
		padding: 12px;
		background-color: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 16px;
		cursor: pointer;
	}

	.login-form button:hover {
		background-color: #0056b3;
	}

	.chat-container {
		background: white;
		border-radius: 8px;
		box-shadow: 0 2px 10px rgba(0,0,0,0.1);
		overflow: hidden;
	}

	.header {
		background-color: #007bff;
		color: white;
		padding: 15px 20px;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.logout-btn {
		background-color: #dc3545;
		color: white;
		border: none;
		padding: 8px 16px;
		border-radius: 4px;
		cursor: pointer;
	}

	.logout-btn:hover {
		background-color: #c82333;
	}

	.topics {
		background-color: #f8f9fa;
		padding: 15px 20px;
		border-bottom: 1px solid #ddd;
	}

	.topics h3 {
		margin: 0 0 10px 0;
		color: #333;
	}

	.topic-list {
		display: flex;
		gap: 10px;
		flex-wrap: wrap;
	}

	.topic-btn {
		padding: 8px 16px;
		border: 1px solid #007bff;
		background-color: white;
		color: #007bff;
		border-radius: 20px;
		cursor: pointer;
		transition: all 0.2s;
	}

	.topic-btn:hover {
		background-color: #007bff;
		color: white;
	}

	.topic-btn.active {
		background-color: #007bff;
		color: white;
	}

	.messages-container {
		height: 500px;
		display: flex;
		flex-direction: column;
	}

	.messages-header {
		padding: 15px 20px;
		border-bottom: 1px solid #ddd;
		background-color: #f8f9fa;
	}

	.messages-header h3 {
		margin: 0;
		color: #333;
	}

	.messages {
		flex: 1;
		overflow-y: auto;
		padding: 20px;
		background-color: white;
	}

	.message {
		margin-bottom: 15px;
		padding: 12px;
		border-left: 3px solid #007bff;
		background-color: #f8f9fa;
		border-radius: 0 8px 8px 0;
	}

	.message-header {
		display: flex;
		justify-content: space-between;
		margin-bottom: 5px;
	}

	.username {
		font-weight: bold;
		color: #007bff;
	}

	.timestamp {
		color: #666;
		font-size: 0.9em;
	}

	.message-content {
		color: #333;
	}

	.message-input {
		padding: 20px;
		border-top: 1px solid #ddd;
		background-color: #f8f9fa;
	}

	.message-input form {
		display: flex;
		gap: 10px;
	}

	.message-input input {
		flex: 1;
		padding: 12px;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 16px;
	}

	.message-input button {
		padding: 12px 20px;
		background-color: #28a745;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}

	.message-input button:hover {
		background-color: #218838;
	}
</style>