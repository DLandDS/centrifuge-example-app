import type { LoginRequest, LoginResponse, User } from '$lib/types';

const API_BASE_URL = 'http://localhost:3001/api';

class ApiClient {
    private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
        const url = `${API_BASE_URL}${endpoint}`;
        
        const defaultHeaders: HeadersInit = {
            'Content-Type': 'application/json',
        };

        // Add auth token if available
        const token = localStorage.getItem('auth_token');
        if (token) {
            defaultHeaders.Authorization = `Bearer ${token}`;
        }

        const config: RequestInit = {
            headers: { ...defaultHeaders, ...options.headers },
            ...options,
        };

        try {
            const response = await fetch(url, config);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            return await response.json();
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    async login(credentials: LoginRequest): Promise<LoginResponse> {
        return this.request<LoginResponse>('/login', {
            method: 'POST',
            body: JSON.stringify(credentials),
        });
    }

    async getUserInfo(): Promise<User> {
        return this.request<User>('/user');
    }

    async refreshCentrifugeToken(): Promise<{ centrifuge_token: string }> {
        return this.request<{ centrifuge_token: string }>('/centrifuge-token', {
            method: 'POST',
        });
    }

    async healthCheck(): Promise<{ status: string }> {
        return this.request<{ status: string }>('/health');
    }
}

export const apiClient = new ApiClient();