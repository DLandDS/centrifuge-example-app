<script lang="ts">
	import { apiClient } from '$lib/api';
	import { login } from '$lib/stores/auth';
	import type { LoginRequest } from '$lib/types';

	let username = '';
	let password = '';
	let isLoading = false;
	let error = '';

	const handleSubmit = async () => {
		if (!username || !password) {
			error = 'Please enter both username and password';
			return;
		}

		isLoading = true;
		error = '';

		try {
			const credentials: LoginRequest = { username, password };
			const response = await apiClient.login(credentials);
			
			login(response.token, response.centrifuge_token, response.user);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Login failed';
		} finally {
			isLoading = false;
		}
	};
</script>

<div class="login-container">
	<div class="login-form">
		<h1>Chat Login</h1>
		
		<form on:submit|preventDefault={handleSubmit}>
			<div class="form-group">
				<label for="username">Username:</label>
				<input
					id="username"
					type="text"
					bind:value={username}
					placeholder="Enter your username"
					disabled={isLoading}
					required
				/>
			</div>

			<div class="form-group">
				<label for="password">Password:</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					placeholder="Enter your password"
					disabled={isLoading}
					required
				/>
			</div>

			{#if error}
				<div class="error">{error}</div>
			{/if}

			<button type="submit" disabled={isLoading}>
				{isLoading ? 'Logging in...' : 'Login'}
			</button>
		</form>

		<div class="demo-info">
			<p><strong>Demo:</strong> You can use any username and password to login</p>
		</div>
	</div>
</div>

<style>
	.login-container {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		background-color: #f5f5f5;
	}

	.login-form {
		background: white;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
		width: 100%;
		max-width: 400px;
	}

	h1 {
		text-align: center;
		color: #333;
		margin-bottom: 2rem;
	}

	.form-group {
		margin-bottom: 1rem;
	}

	label {
		display: block;
		margin-bottom: 0.5rem;
		color: #555;
		font-weight: 500;
	}

	input {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
		box-sizing: border-box;
	}

	input:focus {
		outline: none;
		border-color: #0066cc;
		box-shadow: 0 0 0 2px rgba(0, 102, 204, 0.2);
	}

	input:disabled {
		background-color: #f5f5f5;
		cursor: not-allowed;
	}

	button {
		width: 100%;
		padding: 0.75rem;
		background-color: #0066cc;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	button:hover:not(:disabled) {
		background-color: #0052a3;
	}

	button:disabled {
		background-color: #ccc;
		cursor: not-allowed;
	}

	.error {
		color: #d32f2f;
		background-color: #ffebee;
		padding: 0.75rem;
		border-radius: 4px;
		margin-bottom: 1rem;
		text-align: center;
	}

	.demo-info {
		margin-top: 1.5rem;
		padding: 1rem;
		background-color: #e3f2fd;
		border-radius: 4px;
		text-align: center;
	}

	.demo-info p {
		margin: 0;
		color: #1976d2;
		font-size: 0.9rem;
	}
</style>